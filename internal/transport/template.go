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
                uuid: document.getElementById('passengerId').value
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
                if (!response.ok) {
                    throw new Error('Network response was not ok: ' + response.statusText);
                }
                return response.json();
            })
            .then(data => {
                alert('Success: ' + JSON.stringify(data, null, 2));
            })
            .catch(error => {
                alert('Error: ' + error.message);
            });
        });
    </script>
</body>
</html>
`
