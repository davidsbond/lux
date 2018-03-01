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

  // Register some middleware
  router.Middleware(middlewareFunc)

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

## handlers

Defining a handler is fairly straightforward. You can have one handler per HTTP method. the signature for any handler function is as follows:

```go
func handler(r lux.Request) (lux.Response, error){
  // handle
}
```

Then you can register your handler function

```go
router.Handler("GET", handler)
```

Whenever a GET request is made to the lambda function, the handler will be called. You can also specify which headers must be present for a request to reach your handler. For example, if we only want JSON requests, we can specify:

```go
router.Handler("GET", handler).Headers("Content-Type", "application/json")
```

In the scenario where a request is made with invalid headers, a `400` response is returned.

## recovery

In the event a process in your handler causes a panic, the router will automatically recover for you. However, if you want to handle recovery yourself, you can provide a custom panic handler. The signature for a panic handler is as follows:

```go
func onPanic(r lux.Request, err error) {
  // handle
}
```

Then you can tell the router to use your custom panic handler:

```go
router.Recovery(onPanic)
```

## logging

The router uses [logrus](https://github.com/sirupsen/logrus), a structured logger. You can either choose to disable the logs of the router or you can provide some configuration for it. AWS automatically logs the output of `stderr` and `stdout`, so you can specify that the router should log to either of these like this:

```go
router.Logging(os.Stdout, &logrus.JSONFormatter{})
```

The second parmeter is a logrus formatter, which will output the logs as JSON. You can also provide a custom formatter, see [logrus' godoc page](https://godoc.org/github.com/sirupsen/logrus#Formatter) for more info on custom formatters

## middleware

You can also provide custom middleware functions that can modify a request before it reaches your handler. The method signature for a middleware function is as follows:

```go
func middleware(r *lux.Request) error {
  // do something to the request
}
```

You can register the middleware like this:

```go
router.Middleware(middleware)
```

If the middleware method returns an error, the generated response will be a `500`. Middleware methods are executed in the order they are registered.