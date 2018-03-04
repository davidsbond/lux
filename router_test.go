package lux_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/davidsbond/lux"
	"github.com/stretchr/testify/assert"
)

func TestLux_NewResponse(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Body          interface{}
		Status        int
		ExpectedError string
	}{
		{Body: "hello world", Status: http.StatusOK},
		{Body: make(chan bool), Status: http.StatusOK, ExpectedError: "failed to encode response body"},
	}

	for _, tc := range tt {
		resp, err := lux.NewResponse(tc.Body, tc.Status)

		if err != nil {
			assert.Contains(t, err.Error(), tc.ExpectedError)
			continue
		}

		assert.Equal(t, tc.Status, resp.StatusCode)
		assert.NotEmpty(t, resp.Body)
	}
}

func TestRouter_UsesMiddleware(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Middleware     lux.HandlerFunc
		Request        lux.Request
		Handlers       map[string]lux.HandlerFunc
		ExpectedBody   string
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
			ExpectedBody:   "\"hello test\"\n",
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
			ExpectedBody:   "\"error\"",
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
		assert.Equal(t, tc.ExpectedBody, resp.Body)
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
			ExpectedStatus: http.StatusNotAcceptable,
			ExpectedError:  "not acceptable",
		},
		// Scenario 3: Handler does not exist
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{},
			ExpectedStatus: http.StatusMethodNotAllowed,
			ExpectedError:  "not allowed",
		},
		// Scenario 4: Invalid GET request with no headers.
		{
			Request: lux.Request{
				HTTPMethod: "GET",
				Headers:    map[string]string{},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusNotAcceptable,
			ExpectedError:  "not acceptable",
		},
		// Scenario 5: Valid DELETE request with only a GET handler registered.
		{
			Request: lux.Request{
				HTTPMethod: "DELETE",
				Headers:    map[string]string{"content-type": "application/json"},
			},
			Handlers:       map[string]lux.HandlerFunc{"GET": getHandler},
			ExpectedStatus: http.StatusMethodNotAllowed,
			ExpectedError:  "not allowed",
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

func getHandler(w lux.ResponseWriter, r *lux.Request) {
	w.Headers().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)

	if err := encoder.Encode("hello test"); err != nil {
		//
	}

}

func recoverHandler(req lux.Request, err error) {
	// Do nothing
}

func panicHandler(w lux.ResponseWriter, r *lux.Request) {
	panic("uh oh")
}

func errorMiddleware(w lux.ResponseWriter, r *lux.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("\"error\""))
}

func middleware(w lux.ResponseWriter, r *lux.Request) {

}
