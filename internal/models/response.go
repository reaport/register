package models

import (
	"errors"
	"time"
)

var (
	ErrTicketNotFound = errors.New("❌ ticket not found")
	ErrFlightNotFound = errors.New("❌ flight not found")
	ErrInternalServer = errors.New("❌internal server error")
	ErrBaggageSize    = errors.New("❌ exceeding the allowed baggage size")
)

func GetCode(message string) int {
	switch message {
	case ErrTicketNotFound.Error():
		return 404
	case ErrInternalServer.Error():
		return 500
	default:
		return 500
	}
}

type PassengerResponse struct {
	FlightName        string    `json:"flightName"`
	DepartureTime     time.Time `json:"departureTime"`
	StartPlantingTime time.Time `json:"startPlantingTime"`
	Seat              string    `json:"seat"`
}
