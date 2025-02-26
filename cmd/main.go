package main

import (
	"github.com/reaport/register/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,                  // Включаем цвета даже если вывод не в TTY (опционально)
		FullTimestamp:   true,                  // Показываем полные временные метки
		TimestampFormat: "2006-01-02 15:04:05", // Формат времени
		DisableSorting:  false,                 // Сортировка полей)
	})
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initalization config %s", err.Error())
	}
	app.Run()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
