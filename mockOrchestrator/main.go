package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// finishRegistrationHandler - обработчик для /registration/{flightId}/finish
func finishRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Выводим JSON как строку
	log.Println("Request Body:", string(body))
	// Извлекаем flightId из параметров URL
	//vars := mux.Vars(r)

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Print("Ok", r.Body)
}

func main() {
	// Создаем маршрутизатор
	router := mux.NewRouter()

	// Регистрируем маршрут
	router.HandleFunc("/registration/{flightId}/finish", finishRegistrationHandler).Methods(http.MethodPost)

	// Создаем сервер
	server := &http.Server{
		Addr:    ":8087",
		Handler: router,
	}

	// Запускаем сервер
	fmt.Println("Server starting on http://localhost:8087")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
