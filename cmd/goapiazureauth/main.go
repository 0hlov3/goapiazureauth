package main

import (
	"github.com/0hlov3/goapiazureauth/internal/helpers"
	"github.com/0hlov3/goapiazureauth/internal/logger"
	"github.com/0hlov3/goapiazureauth/internal/models"
	"github.com/0hlov3/goapiazureauth/internal/webserver"
	"go.uber.org/zap"
	"os"
)

func main() {
	logLevel := os.Getenv("GO_AZURE_AUTH_LOGLEVEL")
	log := logger.InitializeZapCustomLogger(logLevel)

	config := models.ApiConfig{
		Log: log,
		Azure: models.Azure{
			TenantID: os.Getenv("AZURE_TEST_API_TENANTID"),
			Scope:    os.Getenv("AZURE_TEST_API_AUD"),
		},
	}
	if !helpers.ContainsEmpty(config.Azure.TenantID, config.Azure.Scope) {
		log.Fatal("TenantID or Scope in Variables not set.")
	}

	wconfig := webserver.NewApiConfig(&config)

	s := webserver.NewServer("127.0.0.1", "8081", wconfig)
	if err := s.Server.ListenAndServe(); err != nil {
		config.Log.Fatal("Error starting webserver", zap.Error(err))
	}
}
