package models

type Passenger struct {
	Uuid          string  `json:"passengerId"`
	BaggageWeight float64 `json:"baggageWeight"`
	MealOption    string  `json:"mealOption"`
	SeatClass     string  `json:"seatClass"`
	Have          bool    // true - зарегестрировался
}
