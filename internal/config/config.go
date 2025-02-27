package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MealOption       []string `json:"mealOption"`
	SeatClass        []string `json:"seatClass"`
	MaxBaggage       float64  `json:"maxBaggage"`
	UrlTicketService string   `json:"urlTicketService"`
	UrlOrchestrator  string   `json:"urlOrchestrator"`
}

func NewConfig() (Config, error) {
	// Открываем файл config.json
	file, err := os.Open("internal/config/config.json")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Создаем пустую структуру Config
	var config Config

	// Декодируем JSON из файла в структуру
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
