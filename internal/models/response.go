package models

import (
	"errors"
	"time"
)

var (
	ErrTicketNotFound = errors.New("❌ ticket not found")
	ErrInternalServer = errors.New("❌internal server error")
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
