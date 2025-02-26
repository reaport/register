package transport

import (
	"encoding/json"
	"fmt"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (api *API) RegisterPassenger(w http.ResponseWriter, r *http.Request) {
	logrus.Info("api RegisterPassenger handlers")
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
	resp, err := api.service.RegisterPassenger(passenger)
	if err != nil {
		writeResponse(w, err.Error(), models.GetCode(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)

}

func (api *API) RegisterFlights(w http.ResponseWriter, r *http.Request) {
	logrus.Info("api RegisterFlights handlers")
	var flight models.Flight

	err := json.NewDecoder(r.Body).Decode(&flight)
	if err != nil {
		logrus.Error("❌Api.RegisterPassengerFlight validation error flight 1")
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Проверка валидации
	if flight.FlightId == "" || flight.FlightName == "" || flight.EndRegisterTime.IsZero() ||
		flight.EndRegisterTime.IsZero() || flight.DepartureTime.IsZero() || flight.StartPlantingTime.IsZero() ||
		flight.SeatsAircraft == nil || len(flight.SeatsAircraft) == 0 {
		logrus.Error("❌Api.RegisterPassengerFlight validation error flight 2")
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}
	for _, v := range flight.SeatsAircraft {
		if v.SeatNumber == "" || v.SeatClass == "" {
			logrus.Error("❌Api.RegisterPassengerFlight validation error flight")
			writeResponse(w, ErrValidation, http.StatusBadRequest)
			return
		}
	}
	logrus.Info("✅ Api.RegisterPassengerFlight make get request for get passenger")
	// Формируем URL для GET-запроса
	url := fmt.Sprintf("http://localhost:8086/flight/%s/passengers", flight.FlightId)
	var resp *http.Response
	// Выполняем GET-запрос
	resp, err = http.Get(url)
	if err != nil {
		logrus.Error("❌Api.RegisterPassengerFlight get request break")
		writeResponse(w, ErrInternal, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус-код (что рейс найден и можно получать пассажиров)
	if resp.StatusCode != http.StatusOK {
		writeResponse(w, ErrInternal, http.StatusInternalServerError)
		return
	}

	// Читаем и парсим тело ответа
	var passengers []models.Passenger
	err = json.NewDecoder(resp.Body).Decode(&passengers)
	if err != nil {
		logrus.Error("❌Api.RegisterPassengerFlight validation error passengers")
		writeResponse(w, ErrValidation, http.StatusBadRequest)
		return
	}

	if len(passengers) > len(flight.SeatsAircraft) {
		logrus.Error("❌Api.RegisterPassengerFlight 👤 unexpected overbooking: ", " flight: ", flight.FlightId)
		writeResponse(w, ErrInternal, http.StatusInternalServerError)
		return
	}
	logrus.Info("✅ Api.RegisterPassengerFlight make register flight")

	err = api.service.RegisterFlights(flight, passengers)
	if err != nil {
		writeResponse(w, err.Error(), models.GetCode(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) GetData(w http.ResponseWriter, r *http.Request) {
	data := api.service.GetData()
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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
