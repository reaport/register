package models

type Passenger struct {
	Uuid          string  `json:"idTraveler"`
	BaggageWeight float64 `json:"baggageWeight"`
	MealOption    string  `json:"mealOption"`
	SeatClass     string  `json:"seatClass"`
}
