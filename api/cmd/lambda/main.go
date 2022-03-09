package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pbivrell/lights/api/app"
	handler "github.com/pbivrell/lights/api/handlers/lambda"
	"github.com/pbivrell/lights/api/storage/dynamo"
)

func main() {

	app := app.New(dynamo.New())

	router := handler.GetRouter(app)

	lambda.Start(router)
}
