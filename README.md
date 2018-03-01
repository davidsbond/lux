# lux

[![CircleCI](https://img.shields.io/circleci/project/github/davidsbond/lux.svg)](https://circleci.com/gh/davidsbond/lux)
[![GoDoc](https://godoc.org/github.com/davidsbond/lux?status.svg)](http://godoc.org/github.com/davidsbond/lux)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidsbond/lux)](https://goreportcard.com/report/github.com/davidsbond/lux)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/davidsbond/lux/release/LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdavidsbond%2Flux.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdavidsbond%2Flux?ref=badge_shield)

An HTTP router for Golang lambda functions

## usage

```go
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

  // Start the lambda.
  lambda.Start(router.HandleRequest)
}
```