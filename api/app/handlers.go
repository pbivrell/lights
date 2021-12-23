package app

import (
	"github.com/gorilla/mux"
	"github.com/pbivrell/simpleserver"
)

type App struct {
	s *simpleserver.Server
}

func NewApp() *App {
	return &App{
		s: simpleserver.NewServer(
			simpleserver.WithHandler(registerHandlers())
		),
	}

}

func registerHandlers() http.Handler {
	r := mux.Router()
	r.HandleFunc("/lights", SetLights).Methods("POST")
	r.HandleFunc("/lights", GetLights).Methods("GET")
	return r
}

func (a *App) SetLights(w http.ResponseWritter, r *http.Request) {

}

func (a *App) GetLights(w http.ResponseWriter, r *http.Request) {

}
