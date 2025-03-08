package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MealOption           []string `json:"mealOption"`
	MaxBaggage           float64  `json:"maxBaggage"`
	UrlTicketService     string   `json:"urlTicketService"`
	UrlOrchestrator      string   `json:"urlOrchestrator"`
	ProdUrlTicketService string
	ProdUrlOrchestrator  string
	MockUrlTicketService string `json:"mockUrlTicketService"`
	MockUrlOrchestrator  string `json:"mockUrlOrchestrator"`
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
	config.ProdUrlOrchestrator = config.UrlOrchestrator
	config.ProdUrlTicketService = config.UrlTicketService
	return config, nil
}
