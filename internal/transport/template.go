package transport

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
            margin: 0;
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
        /* Стили для select */
    select {
        width: 100%;
        padding: 10px;
        margin-bottom: 15px;
        border: 1px solid #dcdcdc;
        border-radius: 5px;
        font-size: 14px;
        background-color: #fff;
        color: #34495e;
        box-sizing: border-box;
        appearance: none; /* Убираем стандартную стрелку */
        background-image: url('data:image/svg+xml;utf8,<svg fill="%2334495e" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg"><path d="M7 10l5 5 5-5z"/></svg>');
        background-repeat: no-repeat;
        background-position: right 10px center;
        background-size: 12px;
        cursor: pointer;
        transition: border-color 0.3s ease, box-shadow 0.3s ease;
    }
    select:focus {
        border-color: #3498db;
        outline: none;
        box-shadow: 0 0 5px rgba(52, 152, 219, 0.3);
    }
    select:hover {
        border-color: #bdc3c7;
    }
    select option {
        padding: 10px;
        background-color: #fff;
        color: #34495e;
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
            margin-bottom: 10px; /* Добавляем отступ между кнопками */
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
            <h5>Now : {{.UrlTicketService}}</h5>
            <select id="urlTicketService" name="urlTicketService">
            <option value="{{.ProdUrlTicketService}}">🎟 Ticket Service</option>
            <option value="{{.MockUrlTicketService}}">🛠 Mock Ticket Service</option>
             </select>

            <label for="urlOrchestrator">Orchestrator Service URL:</label>
            <h5>Now : {{.UrlOrchestrator}}</h5>
            <select id="urlOrchestrator" name="urlOrchestrator">
            <option value="{{.ProdUrlOrchestrator}}">🌍 Orchestrator Service</option>
            <option value="{{.MockUrlOrchestrator}}">🛠Mock Orchestrator Service</option>
             </select>

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

    <!-- Секция с двумя кнопками -->
    <div class="section">
        <h1>📥 Download Data 📥</h1>
        <form action="/download/logs" method="GET">
            <button type="submit">Get logs</button>
        </form>
        <form action="/download/backup" method="GET">
            <button type="submit">Get backup ( info for orchestrator)</button>
        </form>
    </div>

    <script>
        document.getElementById('passengerForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = {
                passengerId: document.getElementById('passengerId').value
            };
            const baggageWeightInput = document.getElementById('baggageWeight').value;
            if (baggageWeightInput) {
                const baggageWeight = parseFloat(baggageWeightInput);
                if (!isNaN(baggageWeight)) {
                    formData.baggageWeight = baggageWeight;
                } else {
                    alert('Error: Baggage Weight must be a valid number');
                    return;
                }
            }
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
            // Парсим JSON независимо от статуса ответа
            return response.json().then(data => {
                if (!response.ok) {
                    // Если статус не 2xx, выбрасываем ошибку с данными из JSON
                    throw new Error(data.errors || 'Unknown error occurred');
                }
                return data; // Успешный ответ
            });
        })
        .then(data => {
            // Успешный случай
            alert('Success: ' + JSON.stringify(data, null, 2));
        })
        .catch(error => {
            // Обработка всех ошибок (сетевых или из JSON)
            alert('Error: ' + error.message);
        });
        });
    </script>
</body>
</html>
`
