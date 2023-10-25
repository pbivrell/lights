package main

import (
	"syscall"

	"github.com/pbivrell/lights/api/app"
	"github.com/pbivrell/lights/api/handlers"
	"github.com/pbivrell/lights/api/handlers/server"
	"github.com/pbivrell/lights/api/storage"
	"github.com/pbivrell/lights/api/storage/mock"
	"github.com/pbivrell/simpleserver"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const (
	DevOrigin  = "http://localhost:3000"
	ProdOrigin = "https://lights.paulbivrell.com"
	CertFile   = "certs/fullchain.pem"
	KeyFile    = "certs/privkey.pem"
)

var AllowOrigins = []string{DevOrigin, ProdOrigin}

func InitDemoStorage(store storage.AppStorage) {

	const (
		DeviceKey     = "theoriginallightskey"
		DeviceID      = "Hub1"
		LightDeviceID = "OriginalLight1"
		ID            = "demo"
	)

	password, _ := app.HashPassword("demo")

	store.User.Write(ID, &storage.User{
		ID:    ID,
		Hash:  password,
		Email: "",
		Hubs: []storage.HubAlias{
			{
				ID:    DeviceID,
				Alias: "",
			},
		},
	})

	store.Hub.Write(DeviceID, &storage.Hub{
		Key:    DeviceKey,
		ID:     DeviceID,
		Lights: []string{},
	})
}

func main() {
	l := logrus.New()
	//l.SetReportCaller(true)
	l.SetLevel(logrus.DebugLevel)

	logger := l.WithFields(logrus.Fields{
		"app": "light-api",
	})
	/*
		store, err := dynamo.NewFromProfile(dynamo.CredsProfile)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("unable to create aws creds")
			return
		}
	*/

	store := mock.AppStorage()

	InitDemoStorage(store)

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
		simpleserver.WithTLS(CertFile, KeyFile),
	)

	go func() {
		s.WithSigShutdown(syscall.SIGTERM)
	}()

	logger.Info("starting server")
	err := s.Run()
	logger.WithFields(logrus.Fields{
		"err": err,
	}).Info("shutting down")

}
