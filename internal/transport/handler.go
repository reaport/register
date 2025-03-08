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
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f4f7fa;
            color: #333;
            margin: 0;
            padding: 20px;
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            font-size: 28px;
            margin-bottom: 30px;
            text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.1);
        }
        h3 {
            color: #2980b9;
            font-size: 20px;
            margin-top: 20px;
        }
        h5 {
            color: #7f8c8d;
            font-size: 16px;
            margin-bottom: 10px;
        }
        .section {
            background: #fff;
            padding: 25px;
            border-radius: 10px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            max-width: 500px;
            margin: 0 auto 30px;
        }
        form {
            margin: 0; /* Убираем лишние отступы формы внутри section */
        }
        label {
            display: block;
            font-size: 14px;
            color: #34495e;
            margin-bottom: 5px;
            font-weight: bold;
        }
        input[type="text"] {
            width: 100%;
            padding: 10px;
            margin-bottom: 15px;
            border: 1px solid #dcdcdc;
            border-radius: 5px;
            font-size: 14px;
            box-sizing: border-box;
            transition: border-color 0.3s ease;
        }
        input[type="text"]:focus {
            border-color: #3498db;
            outline: none;
            box-shadow: 0 0 5px rgba(52, 152, 219, 0.3);
        }
        button {
            background-color: #3498db;
            color: white;
            padding: 12px 20px;
            border: none;
            border-radius: 5px;
            font-size: 16px;
            cursor: pointer;
            width: 100%;
            transition: background-color 0.3s ease;
        }
        button:hover {
            background-color: #2980b9;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            background: #ecf0f1;
            padding: 10px;
            margin-bottom: 5px;
            border-radius: 5px;
            font-size: 14px;
            color: #2c3e50;
        }
    </style>
</head>
<body>
    <div class="section">
        <h1>🚧🛠️ Update URL Configurations 🛠🚧️</h1>
        <form action="/data" method="POST">
            <label for="urlTicketService">Ticket Service URL:</label>
            <input type="text" id="urlTicketService" name="urlTicketService" value="{{.UrlTicketService}}">
            
            <label for="urlOrchestrator">Orchestrator URL:</label>
            <input type="text" id="urlOrchestrator" name="urlOrchestrator" value="{{.UrlOrchestrator}}">

            <label for="maxBaggage">Max Baggage:</label>
            <input type="text" id="maxBaggage" name="maxBaggage" value="{{.MaxBaggage}}">

            <button type="submit">Update</button>
        </form>
    </div>

    <div class="section">
        <h1>✅✈️ Open Flight ✈️✅</h1>
        {{range $flightID, $passengers := .Flights}}
            <h3>Flight ID: {{$flightID}}</h3>
            <h5>Passenger base:</h5>
            <ul>
            {{range $passengers}}
                <li>{{.}}</li>
            {{end}}
            </ul>
        {{end}}
    </div>

    <div class="section">
        <h1>🌏👤 Manual Registration 👤🌏</h1>
        <form id="passengerForm" action="/passenger" method="POST">
            <label for="passengerId">Passenger ID:</label>
            <input type="text" id="passengerId" name="passengerId">
            
            <label for="baggageWeight">Baggage Weight:</label>
            <input type="text" id="baggageWeight" name="baggageWeight">

            <label for="mealOption">Meal Option:</label>
            <input type="text" id="mealOption" name="mealOption">

            <button type="submit">Update</button>
        </form>
    </div>

    <script>
        document.getElementById('passengerForm').addEventListener('submit', function(e) {
            e.preventDefault(); // Отменяем стандартную отправку формы

            const formData = {
                uuid: document.getElementById('passengerId').value // Исправлено на uuid для совместимости с сервером
            };

            // Преобразуем baggageWeight в float, добавляем только если заполнено и валидно
            const baggageWeightInput = document.getElementById('baggageWeight').value;
            if (baggageWeightInput) {
                const baggageWeight = parseFloat(baggageWeightInput);
                if (!isNaN(baggageWeight)) {
                    formData.baggageWeight = baggageWeight;
                } else {
                    alert('Error: Baggage Weight must be a valid number');
                    return; // Прерываем отправку, если введено некорректное значение
                }
            }

            // Добавляем mealOption, только если поле заполнено
            const mealOption = document.getElementById('mealOption').value;
            if (mealOption) {
                formData.mealOption = mealOption;
            }

            fetch('/passenger', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok: ' + response.statusText);
                }
                return response.json(); // Парсим JSON из ответа
            })
            .then(data => {
                // Выводим результат в alert
                alert('Success: ' + JSON.stringify(data, null, 2));
            })
            .catch(error => {
                // Выводим ошибку в alert
                alert('Error: ' + error.message);
            });
        });
    </script>
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
			MaxBaggage       float64
		}{
			UrlTicketService: api.service.Cfg.UrlTicketService,
			UrlOrchestrator:  api.service.Cfg.UrlOrchestrator,
			Flights:          dataFlight,
			MaxBaggage:       api.service.Cfg.MaxBaggage,
		}

		// Рендерим шаблон
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	// Если метод не поддерживается
	default: // Если метод не поддерживается
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
