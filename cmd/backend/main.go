package main

import (
	"fmt"
	"github.com/0hlov3/goapiazureauth/internal/helpers"
	"github.com/0hlov3/goapiazureauth/internal/logger"
	"github.com/0hlov3/goapiazureauth/internal/microsoft"
	"github.com/0hlov3/goapiazureauth/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type httpServer struct {
	Server *http.Server
}

type bconfig struct {
	config *models.BackendConfig
}

func NewBackendConfig(config *models.BackendConfig) *bconfig {
	return &bconfig{config: config}
}

func main() {

	logLevel := os.Getenv("AZURE_TEST_API_LOGLEVEL")
	log := logger.InitializeZapCustomLogger(logLevel)
	log.Info("Logger initialized", zap.Any("Level", log.Level()))
	apiUrl := os.Getenv("AZURE_TEST_API_ENDPOINT")
	if !helpers.ContainsEmpty(apiUrl) {
		log.Fatal("apiUrl in Variables not set.")
	}

	config := models.BackendConfig{
		Log:    log,
		ApiUrl: apiUrl,
	}

	NewEntraConfig := microsoft.NewAzureEntraID(&config)
	nconf := NewBackendConfig(&NewEntraConfig)

	s := NewServer("127.0.0.1", "8082", nconf)
	if err := s.Server.ListenAndServe(); err != nil {
		config.Log.Fatal("Error starting webserver", zap.Error(err))
	}
}

func NewServer(host, port string, c *bconfig) *httpServer {
	router := mux.NewRouter()

	router.HandleFunc("/items", c.getItems)

	s := &httpServer{
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      router,
			WriteTimeout: time.Hour,
			ReadTimeout:  time.Hour,
		},
	}

	c.config.Log.Info(fmt.Sprintf("Server is starting on %s", s.Server.Addr))
	return s
}

// callAPI uses the JWT to authenticate against the API
func (b bconfig) getItems(w http.ResponseWriter, _ *http.Request) {
	token, err := microsoft.ConfigNewClient(b.config)
	if err != nil {
		b.config.Log.Fatal(fmt.Sprintf("Error obtaining JWT: %v", err))
	}
	tokenString := token.AccessToken

	b.config.Log.Debug(fmt.Sprintf("Got Token: %s", tokenString))

	client := resty.New()
	response, err := client.R().
		SetHeader("Authorization", "Bearer "+tokenString).
		Get(b.config.ApiUrl)

	if err != nil {
		b.config.Log.Error(fmt.Sprintf("request error: %v", err))
	}

	if response.StatusCode() != http.StatusOK {
		b.config.Log.Error(fmt.Sprintf("failed to call API: %s", response.String()))
	}

	b.config.Log.Info(fmt.Sprintf("API response: %s\n", response.String()))

	w.WriteHeader(http.StatusOK)
	w.Write(response.Body())
}
