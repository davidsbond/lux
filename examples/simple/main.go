package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davidsbond/lux"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a router
	router := lux.NewRouter()

	// Create a custom panic recovery function (optional). This allows you to do things
	// in the event one of your handlers panics.
	router.Recovery(recoverFunc)

	// Configure the logging (optional), anything in stdout or stderr should be
	// logged by AWS.
	router.Logging(os.Stdout, &logrus.JSONFormatter{})

	// Configure your routes for different HTTP methods. You can specify headers that
	// the request must contain to use this route.
	router.Handler("GET", getFunc)
	router.Handler("PUT", putFunc)
	router.Handler("POST", postFunc)
	router.Handler("DELETE", deleteFunc)

	router.Middleware(middleware)

	// Start the lambda.
	lambda.Start(router.HandleRequest)
}

func recoverFunc(r lux.Request, err error) {
	logrus.WithField("request", r).Errorf("recovered from panic, %v", err.Error())
}

func middleware(r *lux.Request) error {
	// perform some actions on the request
	return nil
}

func getFunc(r lux.Request) lux.Response {
	resp, _ := lux.NewResponse("hello GET request", http.StatusOK)

	return resp
}

func postFunc(r lux.Request) lux.Response {
	resp, _ := lux.NewResponse("hello POST request", http.StatusOK)

	return resp
}

func putFunc(r lux.Request) lux.Response {
	resp, _ := lux.NewResponse("hello PUT request", http.StatusOK)

	return resp
}

func deleteFunc(r lux.Request) lux.Response {
	resp, _ := lux.NewResponse("hello DELETE request", http.StatusOK)

	return resp
}
