package models

type RegistrationFinishRequest struct {
	Meal          map[string]int `json:"meal"`
	BaggageWeight float64        `json:"baggageWeight"`
}

//type Meal struct {
//	TypeMeal string `json:"typeMeal"`
//	Count    int    `json:"count"`
//}
