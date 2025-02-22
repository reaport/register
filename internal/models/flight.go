package models

import "time"

type Flight struct {
	FlightId          string    `json:"flightId"`
	FlightName        string    `json:"flightName"`
	EndRegisterTime   time.Time `json:"endRegisterTime"`
	DepartureTime     time.Time `json:"departureTime"`
	StartPlantingTime time.Time `json:"startPlantingTime"`
	Gate              string    `json:"gate"`
	Terminal          string    `json:"terminal"`
	Aircraft          Aircraft  `json:"aircraft"`
}

type Aircraft struct {
	TotalRows        int    `json:"totalRows"`
	TotalSeatsPerRow int    `json:"totalSeatsPerRow"`
	Rows             []Rows `json:"rows"`
}

type Rows struct {
	RowNumber int     `json:"rowNumber"`
	Seats     []Seats `json:"seats"`
}

type Seats struct {
	SeatNumber string `json:"seatNumber"`
	SeatType   string `json:"seatType"`
}
