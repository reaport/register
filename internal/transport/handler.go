package transport

import (
	"encoding/json"
	"fmt"
	"github.com/reaport/register/internal/errors"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (api *API) RegisterPassenger(w http.ResponseWriter, r *http.Request) {
	logrus.Info("api RegisterPassenger handlers")
	var passenger models.Passenger

	err := json.NewDecoder(r.Body).Decode(&passenger)
	if err != nil {
		writeResponse(w, errors.ErrValid)
		return
	}

	defer r.Body.Close()

	if passenger.Uuid == "" {
		writeResponse(w, errors.ErrValid)
		return
	}
	resp, err := api.service.RegisterPassenger(passenger)
	if err != nil {
		writeResponse(w, err)
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
		logrus.Error("❌Api.RegisterPassengerFlight validation errors flight 1")
		writeResponse(w, errors.ErrValid)
		return
	}

	defer r.Body.Close()

	// Проверка валидации
	if flight.FlightId == "" || flight.FlightName == "" || flight.EndRegisterTime.IsZero() ||
		flight.EndRegisterTime.IsZero() || flight.DepartureTime.IsZero() || flight.StartPlantingTime.IsZero() ||
		flight.SeatsAircraft == nil || len(flight.SeatsAircraft) == 0 {
		logrus.Error("❌Api.RegisterPassengerFlight validation errors flight 2")
		writeResponse(w, errors.ErrValid)
		return
	}
	for _, v := range flight.SeatsAircraft {
		if v.SeatNumber == "" || v.SeatClass == "" {
			logrus.Error("❌Api.RegisterPassengerFlight validation errors flight")
			writeResponse(w, errors.ErrValid)
			return
		}
	}
	logrus.Info("✅ Api.RegisterPassengerFlight make get request for get passenger")
	// Формируем URL для GET-запроса
	url := fmt.Sprintf(api.service.Cfg.UrlTicketService, flight.FlightId)
	var resp *http.Response
	// Выполняем GET-запрос
	resp, err = http.Get(url)
	if err != nil {
		logrus.Error("❌Api.RegisterPassengerFlight get request break")
		writeResponse(w, errors.ErrInternalServer)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус-код (что рейс найден и можно получать пассажиров)
	if resp.StatusCode != http.StatusOK {
		writeResponse(w, errors.ErrInternalServer)
		return
	}

	// Читаем и парсим тело ответа
	var passengers []models.Passenger
	err = json.NewDecoder(resp.Body).Decode(&passengers)
	if err != nil {
		logrus.Error("❌Api.RegisterPassengerFlight validation errors passengers")
		writeResponse(w, errors.ErrInternalServer)
		return
	}

	if len(passengers) > len(flight.SeatsAircraft) {
		logrus.Error("❌Api.RegisterPassengerFlight 👤 unexpected overbooking: ", " flight: ", flight.FlightId)
		writeResponse(w, errors.ErrInternalServer)
		return
	}
	logrus.Info("✅ Api.RegisterPassengerFlight make register flight")

	err = api.service.RegisterFlights(flight, passengers)
	if err != nil {
		writeResponse(w, err)
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

func writeResponse(w http.ResponseWriter, err error) {
	errorResponse := ErrorResponse{Message: err.Error()}
	jsonResponse, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errors.GetCode(err.Error()))
	w.Write(jsonResponse)
}
