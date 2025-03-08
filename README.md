# ✈️ Модуль регистрации пассажиров на рейс
[![Swagger](https://img.shields.io/badge/Swagger-Docs-brightgreen?logo=swagger)](https://github.com/reaport/docs/tree/feat/Register)
[![GoogleDocs](https://img.shields.io/badge/GoogleDocs-Docs-blue?logo=googleDocs)](https://docs.google.com/document/d/1-A99pLnf-T3KJgUowspAIestsUUSzbDQ0Sfr5KvSmdI/edit?tab=t.bpkqrrz6nfsl)

Модуль для регистрации пассажиров с выбором питания и сдачей багажа.

---

## 📋 Описание

- **Время регистрации**: Открытие за **M минут**, закрытие за **P минут** до вылета.
- **Условия**: Только для пассажиров с билетом на рейс, доступного для регистрации.
- **Функции**:
    - Смена типа питания.
    - Сдача багажа с учётом веса (есть ограничения по весу).
---

## 🚀 Установка и запуск

Для запуска приложения можно возпользоваться командами:
* ``make run``
* ``go run cmd/main.go``

#### Mock - сервисы :
*  ``make run_ticket`` покупка билетов
*  ``make run_orchestrator`` оркестратор

### ⚙️ Конфигурация 
`Файл: config.json`

```json
{
  "mealOption":  ["Standard", "Vegetarian", "Vegan", "Gluten-Free"],
  "maxBaggage" : 20.0,
  "urlTicketService": "https://tickets.reaport.ru/flight/%s/passengers",
  "urlOrchestrator": "http://localhost:8087/registration/%s/finish",
  "mockUrlTicketService": "http://localhost:8086/flight/%s/passengers",
  "mockUrlOrchestrator": "http://localhost:8087/registration/%s/finish"
}
```

* `mealOption` - типы питания(динамическое изменение)
* `maxBaggage` - Максимально возможный размер багажа
* `urlTicketService` - url модуля покупки билетов
* `urlOrchestrator` - url модуля оркестратора
*  `mockUrlTicketService` - url mock- модуля покупки билетов
* `mockUrlOrchestrator` - url mock -модуля оркестратора


### 🛠Админка
* ⚖️ Настройка максимально допустимого веса багажа.
* 🔗 Изменение URL сервисов (ticket и orchestrator).
* 👤 Ручная регистрация на рейс.
* 📊 Просмотр рейсов.

![img.png](img.png)


