package main

import (
	"syscall"

	"github.com/pbivrell/lights/api/app"
	"github.com/pbivrell/lights/api/handlers"
	"github.com/pbivrell/lights/api/handlers/server"
	"github.com/pbivrell/lights/api/storage/dynamo"
	"github.com/pbivrell/simpleserver"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const (
	DevOrigin = "http://localhost:3000"
)

var AllowOrigins = []string{DevOrigin}

func main() {
	l := logrus.New()
	//l.SetReportCaller(true)
	l.SetLevel(logrus.DebugLevel)

	logger := l.WithFields(logrus.Fields{
		"app": "light-api",
	})
	store, err := dynamo.NewFromProfile(dynamo.CredsProfile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("unable to create aws creds")
		return
	}

	var handlerResponder handlers.Responder = handlers.NewBaseResponder()

	handlerResponder = handlers.NewDebugResponder(handlerResponder, logger)

	a := app.New(store, handlerResponder, logger)

	s := simpleserver.NewServer(
		simpleserver.WithCorsHandler(
			server.RegisterHandlers(a),
			cors.Options{
				AllowedOrigins:   AllowOrigins,
				AllowCredentials: true,
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			},
		),
	)

	go func() {
		s.WithSigShutdown(syscall.SIGTERM)
	}()

	logger.Info("starting server")
	err = s.Run()
	logger.WithFields(logrus.Fields{
		"err": err,
	}).Info("shutting down")

}
