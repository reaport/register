package transport

import (
	"encoding/json"
	"github.com/reaport/register/internal/models"
	"net/http"
)

func (api *API) RegisterPassenger(w http.ResponseWriter, r *http.Request) {
	var passenger models.Passenger

	err := json.NewDecoder(r.Body).Decode(&passenger)
	if err != nil {
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if passenger.Uuid == "" {
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}

	err = api.service.RegisterPassenger(passenger)
	if err != nil {
		// Todo: SWITCH ошибок и resp, status code
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) RegisterFlights(w http.ResponseWriter, r *http.Request) {
	var flight models.Flight

	err := json.NewDecoder(r.Body).Decode(&flight)
	if err != nil {
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if flight.FlightId == "" || flight.FlightName == "" || flight.Gate == "" || flight.Terminal == "" ||
		flight.Aircraft.TotalRows == 0 || flight.Aircraft.TotalSeatsPerRow == 0 || len(flight.Aircraft.Rows) == 0 {
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}

	for _, row := range flight.Aircraft.Rows {
		if row.RowNumber == 0 || len(row.Seats) == 0 {
			writeResponse(w, ErrValidation, http.StatusBadRequest)
			return
		}
		for _, seat := range row.Seats {
			if seat.SeatNumber == "" || seat.SeatClass == "" {
				writeResponse(w, ErrValidation, http.StatusBadRequest)
				return
			}
		}
	}
	// Todo: Проверить совпадание кол-ва мест и переданных данных (никакие опции мест не потерялись)
	passengers := make([]models.Passenger, flight.Aircraft.TotalRows*flight.Aircraft.TotalSeatsPerRow)
	// Todo: проверка на овербукинг (по длине слайса)
	// Todo: сделать запрос Ане и получить пассажиров ( не забыть сравнить значения питания и остального)

	err := api.service.RegisterFlights(flight, passengers)
	if err != nil {
		// Todo: SWITCH ошибок и resp, status code
	}
	// 7. Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	// Запрос Ане для получения пассажира
}

func (api *API) Administer(w http.ResponseWriter, r *http.Request) {

}

func writeResponse(w http.ResponseWriter, message string, code int) {
	errorResponse := ErrorResponse{Message: message}
	jsonResponse, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
