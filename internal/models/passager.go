package models

type Passenger struct {
	Uuid          string  `json:"idTraveler"`
	BaggageWeight float64 `json:"baggageWeight"`
	FoodOnBoard   string  `json:"foodOnBoard"`
}
