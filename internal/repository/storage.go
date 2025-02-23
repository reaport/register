package repository

import (
	"github.com/reaport/register/internal/apperrors"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"sync"
)

type Storage struct {
	flights []flight // рейсы открытые на регистрацию
	mu      sync.Mutex
}

type flight struct {
	flights       models.Flight
	passenger     []models.Passenger
	spaceAircraft map[string]seat // рассадка пассажиров по самолёту (ключ - место)
	mu            sync.Mutex
	//muPassenger   sync.RWMutex
}

type seat struct {
	seatClass string
	employ    bool // true - место занято, false - пусто
}

func NewStorage() *Storage {
	return &Storage{
		flights: make([]flight, 0),
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
			if seat, exists := s.flights[i].spaceAircraft[seatNumber]; exists && seat.employ == false {
				// Место уже заняли
				if seat.employ {
					return apperrors.ErrSeatTaken
				}
				logrus.Info("✅ RegisterPassengerFlight 👤", " flight: ", flightId, " place: ", seatNumber, "taken passenger: ", passenger.Uuid)
				// Занимаем место
				s.flights[i].mu.Lock()
				seat.employ = true
				s.flights[i].mu.Unlock()
				return nil
			} else {
				logrus.Error("❌RegisterPassengerFlight 👤 error: ", apperrors.ErrSeatNotFound, " flight: ", flightId, " place: ", seatNumber, "already taken another passenger")
				return apperrors.ErrSeatNotFound
			}
		}
	}
	logrus.Error("❌RegisterPassengerFlight 👤 error: ", apperrors.ErrFlightNotFound, " flight: ", flightId)
	return apperrors.ErrFlightNotFound
}

// RegisterFlights - сздание нового рейса, открытого на регистрацию и карту самолёта
func (s *Storage) RegisterFlights(fl models.Flight, spaceAircraft map[string]seat) error {
	// Добавляем открытый рейс
	s.mu.Lock()
	logrus.Info("✅ RegisterFlights ✈️ flight: ", fl.FlightId, " ", fl.FlightName)
	// Добавление в хранилку
	s.flights = append(s.flights, flight{
		flights:       fl,
		spaceAircraft: spaceAircraft,
	})
	s.mu.Unlock()
	return nil
}

func (s *Storage) RemoveFlight(flightId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.flights); i++ {
		if s.flights[i].flights.FlightId == flightId {
			logrus.Info("✅ RemoveFlight️ ✈️ 🗑️  flight: ", flightId, " ", s.flights[i].flights.FlightName)
			s.flights = append(s.flights[:i], s.flights[i+1:]...)
			return nil
		}
	}
	logrus.Info("❌ RemoveFlight️ ✈️ 🗑️  flight: ", flightId)
	return apperrors.ErrFlightNotFound
}

func (s *Storage) GetSpaceAircraft(flightId string) (map[string]seat, error) {
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

func (s *Storage) ExistPassengerInFlight(passengerId, flightId string) bool {
	for i := 0; i < len(s.flights); i++ {
		// Ищем рейс
		if s.flights[i].flights.FlightId == flightId {
			for _, human := range s.flights[i].passenger {
				if human.Uuid == passengerId {
					return true
				}
			}
		}
	}
	return false
}

//// Идём по номеру ряда
//for col := 0; col < len(s.flights[i].spaceAircraft); col++ {
//	// Идём по самому ряду (по строчке)
//	for row := 0; row < len(s.flights[i].spaceAircraft[i]); row++ {
//
//	}
//}
