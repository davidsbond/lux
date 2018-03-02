package lux_test

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/davidsbond/lux"
	"github.com/stretchr/testify/assert"
)

func TestRouter_UsesMiddleware(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Middleware     lux.MiddlewareFunc
		Request        lux.Request
		Handlers       map[string]lux.HandlerFunc
		ExpectedError  string
		ExpectedStatus int
	}{
		// Scenario 1: Valid request & happy path middleware
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusOK,
			Middleware:     middleware,
		},
		// Scenario 2: Valid request but middleware returns an error.
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusInternalServerError,
			Middleware:     errorMiddleware,
			ExpectedError:  "\"error\"",
		},
	}

	for _, tc := range tt {
		// GIVEN that we have a router
		router := lux.NewRouter()
		router.Logging(bytes.NewBuffer([]byte{}), &logrus.JSONFormatter{})

		// AND that router has registered handlers
		for method, handler := range tc.Handlers {
			router.Handler(method, handler).Headers("content-type", "application/json")
		}

		// AND the router has registered middleware
		router.Middleware(tc.Middleware)

		// WHEN we perform a request
		resp, _ := router.HandleRequest(tc.Request)

		// THEN the status code & body should be what we expect.
		assert.Equal(t, tc.ExpectedError, resp.Body)
		assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)
	}
}

func TestRouter_HandlesRequests(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Request        lux.Request
		Handlers       map[string]lux.HandlerFunc
		ExpectedError  string
		ExpectedStatus int
	}{
		// Scenario 1: Valid GET request with correct headers.
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusOK,
		},
		// Scenario 2: Invalid GET request with incorrect header value.
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/xml"},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "cannot determine route for request, check your HTTP method & headers are valid",
		},
		// Scenario 3: Handler does not exist
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "cannot determine route for request, check your HTTP method & headers are valid",
		},
		// Scenario 4: Invalid GET request with no headers.
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "cannot determine route for request, check your HTTP method & headers are valid",
		},
		// Scenario 5: Valid DELETE request with only a GET handler registered.
		{
			Request: lux.Request{
				HTTPMethod: "DELETE",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "cannot determine route for request, check your HTTP method & headers are valid",
		},
	}

	for _, tc := range tt {
		// GIVEN that we have a router
		router := lux.NewRouter()
		router.Logging(bytes.NewBuffer([]byte{}), &logrus.JSONFormatter{})

		// AND that router has handlers registered
		for method, handler := range tc.Handlers {
			router.Handler(method, handler).Headers("content-type", "application/json")
		}

		// WHEN we perform the request
		resp, err := router.HandleRequest(tc.Request)

		// THEN any errors should be what we expect
		if err != nil {
			assert.Equal(t, tc.ExpectedError, err.Error())
		}

		// AND the status code should be what we expect.
		assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)
	}
}

func TestRouter_Recovers(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Request        lux.Request
		Handlers       map[string]lux.HandlerFunc
		ExpectedError  string
		ExpectedStatus int
	}{
		// Scenario 1: Handler panics
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers: map[string]lux.HandlerFunc{"GET": panicHandler},
		},
	}

	for _, tc := range tt {
		// GIVEN that we have a router with a recovery handler.
		router := lux.NewRouter().Recovery(recoverHandler)
		router.Logging(bytes.NewBuffer([]byte{}), &logrus.JSONFormatter{})

		// AND that router has handlers registered
		for method, handler := range tc.Handlers {
			router.Handler(method, handler).Headers("content-type", "application/json")
		}

		// WHEN we perform the request that will panic
		resp, err := router.HandleRequest(tc.Request)

		// THEN any errors should be what we expect
		if err != nil {
			assert.Equal(t, tc.ExpectedError, err.Error())
		}

		// AND the status code should be what we expect.
		assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)
	}
}

func getHandler(lux.Request) lux.Response {
	return lux.Response{
		StatusCode: http.StatusOK,
	}
}

func recoverHandler(req lux.Request, err error) {
	// Do nothing
}

func panicHandler(req lux.Request) lux.Response {
	panic("uh oh")
}

func errorMiddleware(req *lux.Request) error {
	return errors.New("error")
}

func middleware(req *lux.Request) error {
	return nil
}
