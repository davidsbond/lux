# lux

[![CircleCI](https://circleci.com/gh/davidsbond/lux/tree/develop.svg?style=shield)](https://circleci.com/gh/davidsbond/lux)
[![Coverage Status](https://coveralls.io/repos/github/davidsbond/lux/badge.svg?branch=develop)](https://coveralls.io/github/davidsbond/lux?branch=develop)
[![GoDoc](https://godoc.org/github.com/davidsbond/lux?status.svg)](http://godoc.org/github.com/davidsbond/lux)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidsbond/lux)](https://goreportcard.com/report/github.com/davidsbond/lux)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/davidsbond/lux/release/LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdavidsbond%2Flux.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdavidsbond%2Flux?ref=badge_shield)

A simple package for creating RESTful AWS Lambda functions in Go. Inspired by packages like [mux](https://github.com/gorilla/mux) & [negroni](https://github.com/urfave/negroni)

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

  // Configure your routes for different HTTP methods. You can specify headers/params that
  // the request must contain to use this route.
  router.Handler("GET", getFunc).Queries("key", "*")
  router.Handler("PUT", putFunc).Headers("Content-Type", "application/json")
  router.Handler("POST", postFunc).Headers("Content-Type", "application/json")
  router.Handler("DELETE", deleteFunc).Queries("key", "*")

  // Start the lambda.
  lambda.Start(router.ServeHTTP)
}
```

## handlers

Defining a handler is fairly straightforward. You can have multiple handlers per HTTP method. This package attempts to make creating HTTP handlers as similar to the standard library as possible, so provides a signature mirroring a standard HTTP handler. The signature for any handler function is as follows:

```go
func handler(w lux.ResponseWriter, r *lux.Request) {
  encoder := json.NewEncoder(w)

  if err := encoder.Encode("hello world"); err != nil {
    // handle
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
}
```

Then you can register your handler function using the `Router.Handler` method.

```go
router.Handler("GET", handler)
```

Whenever a GET request is made to the lambda function, the handler will be called. You can also specify which headers must be present for a request to reach your handler. For example, if we only want JSON requests, we can use the `Route.Headers` method to specify this:

```go
// Require header with value
router.Handler("GET", handler).Headers("Content-Type", "application/json")

// Require header regardless of value
router.Handler("GET", handler).Headers("Content-Type", "*")
```

We can also perform the same route matching based on query parameters that you would typically see in GET/DELETE requests by using the `Router.Queries` method:

```go
// Require the query parameter with value
router.Handler("GET", handler).Queries("key", "value")

// Require query parameter regardless of value
router.Handler("GET", handler).Queries("key", "*")

// You can have multiple routes for a single HTTP method that expect different query parameters
router.Handler("GET", handler1).Queries("id", "*")
router.Handler("GET", handler2).Queries("name", "*")
```

## recovery

In the event a process in your handler causes a panic, the router will automatically recover for you. However, if you want to handle recovery yourself, you can provide a custom panic handler. The signature for a panic handler is as follows:

```go
func onPanic(info lux.PanicInfo) {
  // do something with the panic information
}
```

The `PanicInfo` type contains the error, stack & request regarding the panic. You can tell the router to use your custom panic handler like so:

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

You can also provide custom middleware functions that can are executed before your handler. These can be registered globally or per-route. You can prevent execution of your handler by using `w.WriteHeader` method. Any modifications to the response writer that occur during execution of middleware functions will create a response and prevent execution of the handler. Middleware methods are executed in the order they are registered.

```go
func middleware(w lux.ResponseWriter, r *lux.Request) {
  // use an error status to prevent further execution
  w.WriteHeader(http.StatusInternalServerError)

  // changes you make to the request object are propagated to
  // your handlers
  r.Body = "you've changed"
}
```

You can register the middleware like this:

```go
// Global middleware
router.Middleware(middleware)

// Route specific middleware
router.Handler("GET", getHandler).Middleware(middleware)
```
