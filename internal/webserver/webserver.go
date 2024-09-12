package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/0hlov3/goapiazureauth/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type HttpServer struct {
	Server *http.Server
}

type apiConfig struct {
	config *models.ApiConfig
}

func NewApiConfig(config *models.ApiConfig) *apiConfig {
	return &apiConfig{config: config}
}

func NewServer(host, port string, c *apiConfig) *HttpServer {
	router := mux.NewRouter()

	noAuthRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Header.Get("Authorization") == ""
	}).Subrouter()

	authRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	noAuthRouter.HandleFunc("/status", c.getHealth)
	authRouter.HandleFunc("/items", c.getItems)
	authRouter.Use(c.middleware)

	s := &HttpServer{
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

func (c apiConfig) getItems(w http.ResponseWriter, _ *http.Request) {
	c.config.Log.Info("Received request for item list")

	items := []models.ItemList{
		{Id: uuid.New(), Item: "Sureki Zealot's Insignia", Level: 639},
		{Id: uuid.New(), Item: "Seal of the Poisoned Pact", Level: 639},
		{Id: uuid.New(), Item: "Spymaster's Web", Level: 639},
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (c apiConfig) getHealth(w http.ResponseWriter, _ *http.Request) {
	c.config.Log.Info("Received health check request")

	health := models.Health{Status: models.StatusOk}
	respondWithJSON(w, http.StatusOK, health)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}
