package main

import (
	"go-backend/internal"
	"go-backend/internal/config"
	"go-backend/internal/handler"
	"go-backend/internal/service"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	conf, err := config.LoadConfig()
	conf.Logger, _ = zap.NewProduction()

	if err != nil {
		conf.Logger.Warn("Configuration loading error: " + err.Error())
		conf.Logger.Info("Using default configuration.")
	} else {
		conf.Logger.Info("Config loaded successfully.")
	}

	conf.Logger.Info("Starting the application for environment: " + conf.Environment)

	s := service.New(make(map[string]internal.Entry), make(map[string][]string))
	h := handler.Handler(s, conf)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h,
	}

	srv.ListenAndServe()
}
