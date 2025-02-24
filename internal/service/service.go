package service

import (
	"github.com/reaport/register/internal/models"
	"github.com/reaport/register/internal/repository"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo *repository.Storage
}

func NewService(repo *repository.Storage) *Service {
	return &Service{repo: repo}
}
func (s *Service) RegisterPassenger(passenger models.Passenger) (models.PassengerResponse, error) {

	// Регистрируем место
	passengerResponse, err := s.repo.RegisterPassengerFlight(passenger)
	if err != nil {
		logrus.Info("❗❗❗Service.RegisterPassenger registration not success✈️🪑❌")
		return models.PassengerResponse{}, err
	}
	return passengerResponse, nil
}

func (s *Service) RegisterFlights(flight models.Flight, passengers []models.Passenger) error {
	// Регистрация рейса
	err := s.repo.RegisterFlights(flight, passengers)
	if err != nil {
		logrus.Error("❌✈️ Service.RegisterFlights not success flight: ", flight.FlightId, " error:", err.Error())
		return err
	}

	// Регистрация рейса
	logrus.Info("✅✈️ Service.RegisterFlights success flight: ", flight.FlightId)
	return nil
}

func (s *Service) Administer() {

}

// Todo: Нужен чек времени закрытия регистрации - например канал
