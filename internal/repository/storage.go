package repository

import (
	"fmt"
	"github.com/reaport/register/internal/config"
	"github.com/reaport/register/internal/errors"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Storage struct {
	flights []Flight // рейсы открытые на регистрацию
	cfg     config.Config
	mu      sync.Mutex
}

type Flight struct {
	flights    models.Flight
	passengers []models.Passenger
	mu         sync.Mutex
}

func NewStorage() *Storage {
	logrus.Info("storage instance initialized")
	return &Storage{
		flights: make([]Flight, 0),
	}
}

func (s *Storage) GetData() map[string][]string {
	logrus.Info("GetData ", s.flights)
	logrus.Info("✅✅✅ GetData ✅✅✅ ")
	data := make(map[string][]string)
	for _, f := range s.flights {
		fmt.Println("✈️ flightId: ", f.flights.FlightId, "✈️ flightName: ", f.flights.FlightName, " seat:", f.flights.SeatsAircraft)
		fmt.Println("End Register Time ", f.flights.EndRegisterTime)
		fmt.Println("👤 passengers", f.passengers)
		for _, pass := range f.passengers {
			data[f.flights.FlightId] = append(data[f.flights.FlightId], pass.Uuid)
		}
	}
	return data
}

// RegisterPassengerFlight - регистрация конкретного пассажира на рейс
func (s *Storage) RegisterPassengerFlight(passenger models.Passenger) (models.PassengerResponse, error) {
	logrus.Info("Storage RegisterPassengerFlight")
	//  Получаем рейс и меняем питание
	flightId, humanId, err := s.getFlightAndIndexHumanAndSetMeal(passenger)
	if err != nil {
		return models.PassengerResponse{}, err
	}
	for i := 0; i < len(s.flights); i++ {
		// Если нашли нужный рейс
		if s.flights[i].flights.FlightId == flightId {
			for seatIndex, seat := range s.flights[i].flights.SeatsAircraft {
				// Проверяем соответсвует ли класс и свободно ли место
				if seat.SeatClass == s.flights[i].passengers[humanId].SeatClass && !seat.Employ {
					logrus.Info("✅ Storage.RegisterPassengerFlight 👤", " flight: ", flightId, " place: ", seat.SeatNumber, "taken passenger: ", passenger.Uuid)
					s.flights[i].flights.SeatsAircraft[seatIndex].Employ = true
					s.flights[i].passengers[humanId].Have = true
					return models.PassengerResponse{FlightName: s.flights[i].flights.FlightName, DepartureTime: s.flights[i].flights.DepartureTime, StartPlantingTime: s.flights[i].flights.StartPlantingTime, Seat: seat.SeatNumber}, nil
				}
			}
			logrus.Error("❌Storage.RegisterPassengerFlight 👤 unexpected overbooking: ", " flight: ", flightId, " места ", s.flights[i].flights.SeatsAircraft, "\n passanger", s.flights[i].passengers)

			return models.PassengerResponse{}, errors.ErrInternalServer
		}
	}
	logrus.Error("❌Storage.RegisterPassengerFlight 👤 : ", " flight: ", flightId, "  not found")
	return models.PassengerResponse{}, errors.ErrTicketNotFound
}

// RegisterFlights - создание нового рейса, открытого на регистрацию и карту самолёта
func (s *Storage) RegisterFlights(fl models.Flight, passengers []models.Passenger) (time.Time, error) {
	// Добавляем открытый рейс

	logrus.Info("✅ Storage.RegisterFlights ✈️ register flight: ", fl.FlightId, " ", fl.FlightName)
	// Добавление в хранилку
	s.flights = append(s.flights, Flight{
		flights:    fl,
		passengers: passengers,
	})
	return fl.EndRegisterTime, nil
}

func (s *Storage) RemoveFlight(flightId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			logrus.Info("✅ Storage.RemoveFlight️ ✈️ 🗑  flight: ", flightId, " ", s.flights[i].flights.FlightName)
			s.flights = append(s.flights[:i], s.flights[i+1:]...)
			return nil
		}
	}
	logrus.Info("❌ Storage.RemoveFlight️ ✈️ 🗑  flight: ", flightId)
	return errors.ErrInternalServer
}

// Получение рейса для пассажира
func (s *Storage) getFlightAndIndexHumanAndSetMeal(human models.Passenger) (string, int, error) {
	// Итерация по рейсам
	fmt.Println(human)
	for i := 0; i < len(s.flights); i++ {
		// Итерация по пассажирам этого рейса
		for passengerIndex, passenger := range s.flights[i].passengers {
			if passenger.Uuid == human.Uuid {
				// Меняем питание если появились новые предпочтения
				if human.MealType != "" {
					s.flights[i].passengers[passengerIndex].MealType = human.MealType
				}
				if human.BaggageWeight > 0.0 {
					s.flights[i].passengers[passengerIndex].BaggageWeight = human.BaggageWeight
				}
				return s.flights[i].flights.FlightId, passengerIndex, nil
			}
		}
	}
	return "", 0, errors.ErrTicketNotFound
}

func (s *Storage) GetMealsAndBaggage(flightId string) (models.RegistrationFinishRequest, error) {
	result := models.RegistrationFinishRequest{Meal: make(map[string]int)}
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			for _, human := range s.flights[i].passengers {
				// Если пассажир зарегался
				if human.Have {
					result.BaggageWeight += human.BaggageWeight
					result.Meal[human.MealType] += 1
				}
			}
			return result, nil
		}
	}
	return models.RegistrationFinishRequest{}, errors.ErrInternalServer
}
