package transport

import (
	"github.com/gorilla/mux"
	"github.com/reaport/register/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	RegisterPassenger = "/passenger"
	RegisterFlight    = "/flights"
	Data              = "/data"
	Download          = "/download"
)

type API struct {
	service *service.Service
	router  *mux.Router
}

func NewAPI(service *service.Service) *API {
	logrus.Info("api instance initialized")
	return &API{service: service, router: mux.NewRouter()}
}

func (api *API) Register() {
	logrus.Info("api Register handlers")
	api.router.HandleFunc(RegisterPassenger, api.RegisterPassenger).Methods(http.MethodPost)
	api.router.HandleFunc(RegisterFlight, api.RegisterFlights).Methods(http.MethodPost)
	api.router.HandleFunc(Data, api.DataHandler)
	api.router.PathPrefix(Download).HandlerFunc(api.DownloadHandler).Methods(http.MethodGet)
}
