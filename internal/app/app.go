package app

import (
	"context"
	"github.com/reaport/register/internal/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	srv := new(transport.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server %s", err.Error())
		}
	}()

	logrus.Print("todo server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured while shutting down server %s", err.Error())
	}
}
