package main

import (
	"github.com/reaport/register/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

func main() {
	// Открываем файл для записи логов
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal("Не удалось открыть файл логов: ", err)
	}
	defer logFile.Close()

	// Настраиваем MultiWriter для вывода в консоль и файл
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,                  // Включаем цвета даже если вывод не в TTY (опционально)
		FullTimestamp:   true,                  // Показываем полные временные метки
		TimestampFormat: "2006-01-02 15:04:05", // Формат времени
		DisableSorting:  false,                 // Сортировка полей)
	})
	if err := initConfig(); err != nil {
		logrus.Fatalf("errors initalization config %s", err.Error())
	}
	app.Run()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
