package transport

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	api        *API
	httpServer *http.Server
}

func NewServer(api *API) *Server {
	logrus.Info("server instance initialized")
	return &Server{
		api: api,
	}
}

// Run отвечает за запуск сервера
func (s *Server) Run(port string) error {
	s.api.Register()
	logrus.Info("server run")
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        s.api.router, // маршруты
		MaxHeaderBytes: 1 << 20,      // 1 мб
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
