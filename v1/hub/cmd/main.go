package main

import (
	"auth/data"
	"auth/data/firestore"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	ctx := context.Background()

	db, err := firestore.Connect(ctx)
	if err != nil {
		fmt.Printf("Broken: %v\n", err)
		return
	}
	defer db.Close()

	app := NewApp(db)

	r := mux.NewRouter()
	r.HandleFunc("/con", app.HandleCon)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

type App struct {
	DB *firestore.Connection
}

func NewApp(db *firestore.Connection) *App {
	return &App{
		DB: db,
	}
}

func (a *App) HandleCon(w http.ResponseWriter, r *http.Request) {

	ip := r.URL.Query().Get("ip")
	id := r.URL.Query().Get("id")

	light := data.Light{
		ID:      id,
		IP:      ip,
		Updated: time.Now(),
	}
	err := a.DB.WriteLight(r.Context(), light)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Printf("HandleCon: %v", err)
		return
	}
}
