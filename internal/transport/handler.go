package transport

import (
	"encoding/json"
	"fmt"
	"github.com/reaport/register/internal/errors"
	"github.com/reaport/register/internal/models"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
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

// DataHandler остаётся без изменений
func (api *API) DataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, err := template.New("form").Parse(formTemplate)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		dataFlight := api.service.GetData()
		data := struct {
			UrlTicketService string
			UrlOrchestrator  string
			Flights          map[string][]string
			MaxBaggage       float64
		}{
			UrlTicketService: api.service.Cfg.UrlTicketService,
			UrlOrchestrator:  api.service.Cfg.UrlOrchestrator,
			Flights:          dataFlight,
			MaxBaggage:       api.service.Cfg.MaxBaggage,
		}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case "POST":
		logrus.Info(r.FormValue("urlTicketService"))
		api.service.Cfg.UrlTicketService = r.FormValue("urlTicketService")
		api.service.Cfg.UrlOrchestrator = r.FormValue("urlOrchestrator")
		var err error
		api.service.Cfg.MaxBaggage, err = strconv.ParseFloat(r.FormValue("maxBaggage"), 64)
		if err != nil {
			logrus.Errorf("Error parsing maxBaggage %v", err)
		}
		http.Redirect(w, r, "/data", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Обработчик для скачивания файлов
func (api *API) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("DownloadHandler instance")
	switch r.URL.Path {
	case "/download/logs":
		serveFile(w, r, "app.log", "application/octet-stream", "logs.txt")
	case "/download/backup":
		serveFile(w, r, "backUp.txt", "application/octet-stream", "backup.txt")
	default:
		http.Error(w, "Not Foundd", http.StatusInternalServerError)
	}
}

// Вспомогательная функция для отправки файла
func serveFile(w http.ResponseWriter, r *http.Request, filePath, contentType, downloadName string) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Errorf("Error opening file %s: %v", filePath, err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil {
		logrus.Errorf("Error serving file %s: %v", filePath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeResponse(w http.ResponseWriter, err error) {
	errorResponse := ErrorResponse{Message: err.Error()}
	jsonResponse, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errors.GetCode(err.Error()))
	w.Write(jsonResponse)
}
