package repository

import (
	"github.com/reaport/register/internal/apperrors"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"sync"
)

type Storage struct {
	flights []Flight // рейсы открытые на регистрацию
	mu      sync.Mutex
}

type Flight struct {
	flights       models.Flight
	passengers    []models.Passenger
	spaceAircraft map[string]Seat // рассадка пассажиров по самолёту (ключ - место)
	mu            sync.Mutex
	//muPassenger   sync.RWMutex
}

type Seat struct {
	SeatClass string
	Employ    bool // true - место занято, false - пусто
}

func NewStorage() *Storage {
	return &Storage{
		flights: make([]Flight, 0),
	}
}

// RegisterPassengerFlight - регистрация конкретного пассажира на рейс
func (s *Storage) RegisterPassengerFlight(passenger models.Passenger, flightId string, seatNumber string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.flights); i++ {
		// Если нашли нужный рейс
		if s.flights[i].flights.FlightId == flightId {
			// Проверяем, существует ли место с указанным номером
			if seat, exists := s.flights[i].spaceAircraft[seatNumber]; exists && seat.Employ == false {
				// Место уже заняли
				if seat.Employ {
					return apperrors.ErrSeatTaken
				}
				logrus.Info("✅ Storage.RegisterPassengerFlight 👤", " flight: ", flightId, " place: ", seatNumber, "taken passenger: ", passenger.Uuid)
				// Занимаем место
				s.flights[i].mu.Lock()
				seat.Employ = true
				s.flights[i].mu.Unlock()
				return nil
			} else {
				logrus.Error("❌Storage.RegisterPassengerFlight 👤 error: ", apperrors.ErrSeatNotFound, " flight: ", flightId, " place: ", seatNumber, "already taken another passenger")
				return apperrors.ErrSeatNotFound
			}
		}
	}
	logrus.Error("❌Storage.RegisterPassengerFlight 👤 error: ", apperrors.ErrFlightNotFound, " flight: ", flightId)
	return apperrors.ErrFlightNotFound
}

// RegisterFlights - сздание нового рейса, открытого на регистрацию и карту самолёта
func (s *Storage) RegisterFlights(fl models.Flight, spaceAircraft map[string]Seat, passengers []models.Passenger) error {
	// Добавляем открытый рейс
	s.mu.Lock()
	logrus.Info("✅ Storage.RegisterFlights ✈️ flight: ", fl.FlightId, " ", fl.FlightName)
	// Добавление в хранилку
	s.flights = append(s.flights, Flight{
		flights:       fl,
		spaceAircraft: spaceAircraft,
		passengers:    passengers,
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
	return apperrors.ErrFlightNotFound
}

func (s *Storage) GetSpaceAircraft(flightId string) (map[string]Seat, error) {
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			return s.flights[i].spaceAircraft, nil
		}
	}
	return nil, apperrors.ErrFlightNotFound
}

func (s *Storage) ExistFlight(flightId string) bool {
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			return true
		}
	}
	return false
}

func (s *Storage) GetFlightForPassenger(passengerId string) (string, error) {
	// Итерация по рейсам
	for i := 0; i < len(s.flights); i++ {
		// Итерация по пассажирам этого рейса
		for _, passenger := range s.flights[i].passengers {
			if passenger.Uuid == passengerId {
				return s.flights[i].flights.FlightId, nil
			}
		}
	}
	return "", apperrors.ErrTicketNotFound
}

func (s *Storage) GetSeatForPassenger(flightId string, seatClass string) (string, error) {
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			for seatNumber, seatOption := range s.flights[i].spaceAircraft {
				s.flights[i].mu.Lock()
				// Если место свободно и соответвует классу => занимаем
				if !seatOption.Employ && seatOption.SeatClass == seatClass {
					s.flights[i].spaceAircraft[seatNumber] = Seat{
						SeatClass: seatClass,
						Employ:    true,
					}
					s.flights[i].mu.Unlock()
					return seatNumber, nil
				}
				s.flights[i].mu.Lock()
			}
		}
	}
	return "", apperrors.ErrInternalServer
}

//// Идём по номеру ряда
//for col := 0; col < len(s.flights[i].spaceAircraft); col++ {
//	// Идём по самому ряду (по строчке)
//	for row := 0; row < len(s.flights[i].spaceAircraft[i]); row++ {
//
//	}
//}
