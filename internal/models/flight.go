package models

import "time"

type Flight struct {
	FlightId          string          `json:"flightId"`
	FlightName        string          `json:"flightName"`
	EndRegisterTime   time.Time       `json:"endRegisterTime"`
	DepartureTime     time.Time       `json:"departureTime"`
	StartPlantingTime time.Time       `json:"startPlantingTime"`
	SeatsAircraft     []SeatsAircraft `json:"seatsAircraft"`
}

type SeatsAircraft struct {
	SeatNumber string `json:"seatNumber"`
	SeatClass  string `json:"seatClass"`
	Employ     string `json:"-"`
}
