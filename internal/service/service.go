package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/reaport/register/internal/config"
	"github.com/reaport/register/internal/models"
	"github.com/reaport/register/internal/repository"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Service struct {
	repo *repository.Storage
	cfg  config.Config
}

func NewService(repo *repository.Storage, cfg config.Config) *Service {
	logrus.Info("service instance initialized")
	return &Service{repo: repo, cfg: cfg}
}
func (s *Service) RegisterPassenger(passenger models.Passenger) (models.PassengerResponse, error) {
	logrus.Info("Service RegisterPassenger")
	var i int
	// Проверяем что питание соотв конфигу
	for i = 0; i < len(s.cfg.MealOption); i++ {
		if passenger.MealOption == s.cfg.MealOption[i] {
			break
		}
	}

	if i == len(s.cfg.MealOption) {
		passenger.MealOption = ""
	}

	// Проверяем багаж
	if passenger.BaggageWeight > s.cfg.MaxBaggage {
		return models.PassengerResponse{}, models.ErrBaggageSize
	}
	// Регистрируем место
	passengerResponse, err := s.repo.RegisterPassengerFlight(passenger)
	if err != nil {
		logrus.Info("❗️❗️❗️Service.RegisterPassenger registration not success✈️🪑❌")
		return models.PassengerResponse{}, err
	}
	return passengerResponse, nil
}

func (s *Service) RegisterFlights(flight models.Flight, passengers []models.Passenger) error {
	logrus.Info("Service RegisterFlights")
	// Регистрация рейса
	registationTime, err := s.repo.RegisterFlights(flight, passengers)
	if err != nil {
		logrus.Error("❌✈️ Service.RegisterFlights not success flight: ", flight.FlightId, " error:", err.Error())
		return err
	}
	go s.StopRegister(registationTime, flight.FlightId)
	// Регистрация рейса
	logrus.Info("✅✈️ Service.RegisterFlights success flight: ", flight.FlightId)
	return nil
}

func (s *Service) Administer() {

}

func (s *Service) GetData() []repository.Flight {
	return s.repo.GetData()
}

func (s *Service) StopRegister(registationTime time.Time, flightId string) {
	registationTime = registationTime.Add(-3 * time.Hour)
	if time.Now().After(registationTime) {
		err := s.repo.RemoveFlight(flightId)
		if err != nil {
			logrus.Error("❌✈️ Service.StopRegister not success flight: ", flightId, " error:", err.Error())
		}
		return
	}
	timeUntil := time.Until(registationTime)
	timerChan := time.After(timeUntil)

	select {
	case <-timerChan:
		defer func() {
			err := s.repo.RemoveFlight(flightId)

			if err != nil {
				logrus.Error("❌✈️ Service.StopRegister (timer) not success flight: ", flightId, " error:", err.Error())
			}
		}()
		// Время наступило, отправляем никите и удаляем
		result, err := s.repo.GetMealsAndBaggage(flightId)
		if err != nil {
			logrus.Error("❌✈️ Service.StopRegister not success get info flight ", flightId, " error:", err.Error())
			return
		}
		// ОТправляем Никите
		err = s.SendOrch(result, flightId)
		if err != nil {
			// Если не прошло сохраняем  models.RegistrationFinishRequest в backUp.txt
			logrus.Error("❌✈️ Service.StopRegister not success get send Hikita ", flightId, " error:", err.Error())
		}
	}
}

type RegistrationFinishRequest struct {
	Meal          []Meal  `json:"meal"`
	BaggageWeight float64 `json:"baggageWeight"`
}
type Meal struct {
	TypeMeal string `json:"typeMeal"`
	Count    int    `json:"count"`
}

func (s *Service) SendOrch(reqData models.RegistrationFinishRequest, flightId string) error {
	var mealSlice []Meal
	for typeMeal, count := range reqData.Meal {
		mealSlice = append(mealSlice, Meal{
			TypeMeal: typeMeal,
			Count:    count,
		})
	}

	// Создаем полную структуру для отправки
	req := RegistrationFinishRequest{
		Meal:          mealSlice,
		BaggageWeight: reqData.BaggageWeight,
	}

	// Если время еще не наступило, сохраняем данные в backUp.txt
	file, err := os.Create("backUp.txt")
	if err != nil {
		logrus.Error("❌✈️ Service.StopRegister dont open file backUp ", flightId, " error:", err.Error())
	}
	defer file.Close()
	// Сохраняем данные в JSON
	err = json.NewEncoder(file).Encode("flight_id:" + flightId)
	err = json.NewEncoder(file).Encode(req)
	if err != nil {
		logrus.Error("❌✈️ Service.StopRegister dont save file backUp ", flightId, " error:", err.Error())
	}

	// Логируем преобразованные данные
	logrus.Info("Sending request: ", req)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Формируем URL с flightId
	url := fmt.Sprintf(s.cfg.UrlOrchestrator, flightId)

	// Создаем POST-запрос
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	logrus.Info("✅ Successfully sent request for flight %s\n", flightId)
	return nil
}
