package app

import (
	"context"
	"github.com/reaport/register/internal/config"
	"github.com/reaport/register/internal/repository"
	"github.com/reaport/register/internal/service"
	"github.com/reaport/register/internal/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	repo := repository.NewStorage()
	service := service.NewService(repo, cfg)
	api := transport.NewAPI(service)
	srv := transport.NewServer(api)
	go func() {
		if err := srv.Run(viper.GetString("port")); err != nil {
			logrus.Fatalf("errors occured while running http server %s", err.Error())
		}
	}()

	logrus.Print("register server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("errors occured while shutting down server %s", err.Error())
	}
}
