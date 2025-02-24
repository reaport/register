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
	Cfg     config.Config
	Service *service.Service
}

func NewAPI(service *service.Service, Cfg config.Config) *API {
	return &API{Service: service, Cfg: Cfg}
}

func (api *API) Register(router *mux.Router) {
	router.HandleFunc(RegisterPassenger, api.RegisterPassenger).Methods(http.MethodPost)
	router.HandleFunc(RegisterFlight, api.RegisterFlights).Methods(http.MethodPost)
	router.HandleFunc(AdminParam, api.Administer).Methods(http.MethodPost)
	router.HandleFunc(AdminParam, api.Administer).Methods(http.MethodGet)
}
