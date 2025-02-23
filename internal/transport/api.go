package transport

import (
	"github.com/gorilla/mux"
	"github.com/reaport/register/internal/config"
	"github.com/reaport/register/internal/service"
	"net/http"
)

const (
	RegisterPassenger = "/passenger"
	RegisterFlight    = "/flights"
	AdminParam        = "/admin"
)

type API struct {
	cfg     config.Config
	service service.Service
}

func New(service service.Service) *API {
	return &API{service: service}
}

func (api *API) Register(router *mux.Router) {
	router.HandleFunc(RegisterPassenger, api.RegisterPassenger).Methods(http.MethodPost)
	router.HandleFunc(RegisterFlight, api.RegisterFlights).Methods(http.MethodPost)
	router.HandleFunc(AdminParam, api.Administer).Methods(http.MethodPost)
	router.HandleFunc(AdminParam, api.Administer).Methods(http.MethodGet)
}
