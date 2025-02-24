package repository

import (
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"sync"
)

type Storage struct {
	flights []Flight // рейсы открытые на регистрацию
	mu      sync.Mutex
}

type Flight struct {
	flights    models.Flight
	passengers []models.Passenger
	mu         sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		flights: make([]Flight, 0),
	}
}

// RegisterPassengerFlight - регистрация конкретного пассажира на рейс
func (s *Storage) RegisterPassengerFlight(passenger models.Passenger) (models.PassengerResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//  Получаем рейс и меняем питание
	// TODO Нужна замена питания
	flightId, err := s.getFlightAndSetMeal(passenger.Uuid)
	if err != nil {
		return models.PassengerResponse{}, err
	}
	for i := 0; i < len(s.flights); i++ {
		// Если нашли нужный рейс
		if s.flights[i].flights.FlightId == flightId {
			for seatIndex, seat := range s.flights[i].flights.SeatsAircraft {
				// Проверяем соответсвует ли класс и свободно ли место
				if seat.SeatClass == passenger.SeatClass && !seat.Employ {
					logrus.Info("✅ Storage.RegisterPassengerFlight 👤", " flight: ", flightId, " place: ", seat.SeatNumber, "taken passenger: ", passenger.Uuid)
					s.flights[i].flights.SeatsAircraft[seatIndex].Employ = true
					return models.PassengerResponse{FlightName: s.flights[i].flights.FlightName, DepartureTime: s.flights[i].flights.DepartureTime, StartPlantingTime: s.flights[i].flights.StartPlantingTime, Seat: seat.SeatNumber}, nil
				}
			}
			logrus.Error("❌Storage.RegisterPassengerFlight 👤 unexpected overbooking: ", " flight: ", flightId)
			return models.PassengerResponse{}, models.ErrInternalServer
		}
	}
	logrus.Error("❌Storage.RegisterPassengerFlight 👤 : ", " flight: ", flightId, "  not found")
	return models.PassengerResponse{}, models.ErrInternalServer
}

// RegisterFlights - создание нового рейса, открытого на регистрацию и карту самолёта
func (s *Storage) RegisterFlights(fl models.Flight, passengers []models.Passenger) error {
	// Добавляем открытый рейс
	s.mu.Lock()
	logrus.Info("✅ Storage.RegisterFlights ✈️ register flight: ", fl.FlightId, " ", fl.FlightName)
	// Добавление в хранилку
	s.flights = append(s.flights, Flight{
		flights:    fl,
		passengers: passengers,
	})
	s.mu.Unlock()
	return nil
}

func (s *Storage) RemoveFlight(flightId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			logrus.Info("✅ Storage.RemoveFlight️ ✈️ 🗑️  flight: ", flightId, " ", s.flights[i].flights.FlightName)
			s.flights = append(s.flights[:i], s.flights[i+1:]...)
			return nil
		}
	}
	logrus.Info("❌ Storage.RemoveFlight️ ✈️ 🗑️  flight: ", flightId)
	return models.ErrInternalServer
}

// Получение рейса для пассажира
func (s *Storage) getFlightAndSetMeal(passengerId string) (string, error) {
	// Итерация по рейсам
	for i := 0; i < len(s.flights); i++ {
		// Итерация по пассажирам этого рейса
		for _, passenger := range s.flights[i].passengers {
			if passenger.Uuid == passengerId {
				// TODO Нужна замена питания (сравниваем с конфигом)
				return s.flights[i].flights.FlightId, nil
			}
		}
	}
	return "", models.ErrTicketNotFound
}
