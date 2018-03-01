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
	router.Handler("GET", getFunc).Headers("Content-Type", "application/json")
	router.Handler("PUT", putFunc).Headers("Content-Type", "application/json")
	router.Handler("POST", postFunc).Headers("Content-Type", "application/json")
	router.Handler("DELETE", deleteFunc).Headers("Content-Type", "application/json")

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
	return lux.Response{
		StatusCode: http.StatusOK,
		Body:       "hello GET request",
	}
}

func postFunc(r lux.Request) lux.Response {
	return lux.Response{
		StatusCode: http.StatusOK,
		Body:       "hello POST request",
	}
}

func putFunc(r lux.Request) lux.Response {
	return lux.Response{
		StatusCode: http.StatusOK,
		Body:       "hello PUT request",
	}
}

func deleteFunc(r lux.Request) lux.Response {
	return lux.Response{
		StatusCode: http.StatusOK,
		Body:       "hello DELETE request",
	}
}
