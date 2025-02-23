package service

import (
	"github.com/reaport/register/internal/models"
	"github.com/reaport/register/internal/repository"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo repository.Storage
}

func (s *Service) RegisterPassenger(passenger models.Passenger) error {
	flightId, err := s.repo.GetFlightForPassenger(passenger.Uuid)
	if err != nil {
		logrus.Info("❌🔍  Service.RegisterPassenger 👤 passenger: ", passenger.Uuid, " found is not flight", flightId)
		return err
	}
	logrus.Info("✅🔍  Service.RegisterPassenger 👤 passenger: ", passenger.Uuid, " found is flight", flightId)

	// TODO Нужно замена питания

	// Получаем место для пассажира с нужным классом
	seat, err := s.repo.GetSeatForPassenger(flightId, passenger.SeatClass)
	if err != nil {
		logrus.Info("❗❗❗Service.RegisterPassenger unexpected overbooking️ ✈️🪑❌")
		return err
	}

	// Регистрируем место
	err := s.repo.RegisterPassengerFlight(passenger, flightId, seat)
	if err != nil {
		logrus.Info("❗❗❗Service.RegisterPassenger registration not success✈️🪑❌")
		return err
	}
	return nil
}

func (s *Service) RegisterFlights(flight models.Flight, passengers []models.Passenger) error {
	spaceAircraft := make(map[string]repository.Seat, 0)
	// TODO: Логика преобразования Flight. Aircraft в spaceAircraft map[string]Seat

	// Регистрация рейса
	err := s.repo.RegisterFlights(flight, spaceAircraft, passengers)
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
