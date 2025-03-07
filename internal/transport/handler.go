package transport

import (
	"encoding/json"
	"fmt"
	"github.com/reaport/register/internal/errors"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"html/template"
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
		writeResponse(w, errors.ErrTicketUnavailable)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус-код (что рейс найден и можно получать пассажиров)
	if resp.StatusCode != http.StatusOK {
		writeResponse(w, errors.ErrTicketUnavailable)
		return
	}

	// Читаем и парсим тело ответа
	var passengers []models.Passenger
	err = json.NewDecoder(resp.Body).Decode(&passengers)
	if err != nil {
		logrus.Error("❌Api.RegisterPassengerFlight validation errors passengers")
		writeResponse(w, errors.ErrTicketUnavailable)
		return
	}

	if len(passengers) > len(flight.SeatsAircraft) {
		logrus.Error("❌Api.RegisterPassengerFlight 👤 unexpected overbooking: ", " flight: ", flight.FlightId)
		writeResponse(w, errors.ErrTicketUnavailable)
		return
	}
	logrus.Info("✅ Api.RegisterPassengerFlight make register flight")

	err = api.service.RegisterFlights(flight, passengers)
	if err != nil {
		writeResponse(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

type GetOpenFlightResp struct {
	FlightId   string   `json:"flightId"`
	Passangers []string `json:"passangers"`
}

func (api *API) GetData(w http.ResponseWriter, r *http.Request) {
	data := api.service.GetData()
	jsonResponse, _ := json.Marshal(data)
	fmt.Println("data", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

const formTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Update URLs</title>
</head>
<body>
    <h1>Update URL Configurations</h1>
    <form action="/data" method="POST">
        <label for="urlTicketService">Ticket Service URL:</label>
        <input type="text" id="urlTicketService" name="urlTicketService" value="{{.UrlTicketService}}"><br><br>
        
        <label for="urlOrchestrator">Orchestrator URL:</label>
        <input type="text" id="urlOrchestrator" name="urlOrchestrator" value="{{.UrlOrchestrator}}"><br><br>

        <button type="submit">Update</button>
    </form>

    <h2>Flight Passengers</h2>
    {{range $flightID, $passengers := .Flights}}
        <h3>Flight ID: {{$flightID}}</h3>
        <ul>
        {{range $passengers}}
            <li>{{.}}</li>
        {{end}}
        </ul>
    {{end}}
</body>
</html>
`

// Обработчик для маршрута /data
func (api *API) DataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Отображаем HTML форму с текущими значениями URL
		tmpl, err := template.New("form").Parse(formTemplate)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		dataFlight := api.service.GetData()
		// Формируем данные для шаблона
		data := struct {
			UrlTicketService string
			UrlOrchestrator  string
			Flights          map[string][]string
		}{
			UrlTicketService: api.service.Cfg.UrlTicketService,
			UrlOrchestrator:  api.service.Cfg.UrlOrchestrator,
			Flights:          dataFlight,
		}

		// Рендерим шаблон
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		//case "POST":
		//	// Обновляем значения переменных из формы
		//	urlTicketService = r.FormValue("urlTicketService")
		//	urlOrchestrator = r.FormValue("urlOrchestrator")
		//
		//	// Возвращаем пользователя на страницу с GET запросом
		//	http.Redirect(w, r, "/data", http.StatusSeeOther)
		//default:
		// Если метод не поддерживается
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
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
