package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/networks", Networks)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/ip", IP)

	return r
}

func Networks(w http.ResponseWriter, r *http.Request) {

	networks := []string{
		"RuckenFi",
		"SweetAsWiFi",
		"P-Lights-0",
	}

	json.NewEncoder(w).Encode(networks)

}

func Login(w http.ResponseWriter, r *http.Request) {

	ssid := r.URL.Query().Get("ssid")
	if ssid != "RuckenFi" {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	password := r.URL.Query().Get("password")
	if password != "password" {
		http.Error(w, "", http.StatusUnauthorized)
	}

}

func IP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "localhost:3002")
}
