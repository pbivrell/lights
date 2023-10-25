package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	Cert = "./fullchain.pem"
	Key  = "./privkey.pem"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	/*certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("lights.paulbivrell.com"),
		Cache:      autocert.DirCache("certs"),
	}
	tlsConfig := certManager.TLSConfig()
	*/
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
		//TLSConfig: tlsConfig,
	}

	fmt.Println(server.ListenAndServeTLS(Cert, Key))
}
