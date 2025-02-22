package transport

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	RegisterPassenger = "/passenger"
	RegisterFlight    = "/flights"
	AdminParam        = "/admin"
)

type API struct {
}

func New() *API {
	return &API{}
}

func (api *API) Register(router *mux.Router) {
	router.HandleFunc(RegisterPassenger, api.RegisterPassenger).Methods(http.MethodPost)
	router.HandleFunc(RegisterFlight, api.RegisterFlights).Methods(http.MethodPost)
	router.HandleFunc(AdminParam, api.Administer).Methods(http.MethodPost)
	router.HandleFunc(AdminParam, api.Administer).Methods(http.MethodGet)
}
