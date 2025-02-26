package models

import "time"

type Passenger struct {
	Uuid          string  `json:"passengerId"`
	BaggageWeight float64 `json:"baggageWeight"`
	MealOption    string  `json:"mealOption"`
	SeatClass     string  `json:"seatClass"`
	Have          bool    // true - зарегестрировался
}

type PassengerResponse struct {
	FlightName        string    `json:"flightName"`
	DepartureTime     time.Time `json:"departureTime"`
	StartPlantingTime time.Time `json:"startPlantingTime"`
	Seat              string    `json:"seat"`
}
