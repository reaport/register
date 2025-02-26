package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Структура Passenger для ответа
type Passenger struct {
	PassengerId string `json:"passengerId"`
	MealOption  string `json:"mealOption"`
	SeatClass   string `json:"seatClass"`
}

// Обработчик GET-запроса
func passengersHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это GET-запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем {id} из URL
	path := r.URL.Path
	if !strings.HasPrefix(path, "/flight/") || !strings.Contains(path, "/passengers") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Можно получить flightId, если нужно использовать его в логике
	flightId := strings.TrimPrefix(path, "/flight/")
	flightId = strings.TrimSuffix(flightId, "/passengers")
	// Здесь flightId можно использовать, если нужно, но в данном случае он не влияет на ответ

	// Заранее заданный список пассажиров
	passengers := []Passenger{
		{
			PassengerId: "uuid-1234",
			MealOption:  "да",
			SeatClass:   "business",
		},
		{
			PassengerId: "uuid-5678",
			MealOption:  "нет",
			SeatClass:   "economy",
		},
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Сериализуем данные в JSON и отправляем в ответ
	err := json.NewEncoder(w).Encode(passengers)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Регистрируем обработчик для маршрута
	http.HandleFunc("/flight/", passengersHandler)

	// Запускаем сервер на порту 8086
	fmt.Println("Server started at :8086")
	err := http.ListenAndServe(":8086", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
