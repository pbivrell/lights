package main

import (
	"github.com/pbivrell/simpleserver"
	"github.com/rs/cors"

	"github.com/pbivrell/lights/mock/server"
)

func main() {

	s := simpleserver.NewServer(
		simpleserver.WithPort(3002),
		simpleserver.WithCorsHandler(server.Handlers(), cors.Options{
			AllowedOrigins: []string{"http://localhost:3000"},
		}),
	)

	s.Run()
}
