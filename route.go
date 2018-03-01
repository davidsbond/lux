package lux

import (
	"github.com/aws/aws-lambda-go/events"
)

type (
	// The Route type defines a route that can be used by the router.
	Route struct {
		handler HandlerFunc
		method  string
		headers map[string]string
	}
)

// Headers allows you to specify headers a request should have in order to
// use this route.
func (r *Route) Headers(pairs ...string) *Route {
	// Loop through the headers
	for i := 0; i < len(pairs); i += 2 {
		// If we have an odd number of pairs, skip the last one.
		if len(pairs) < i+1 {
			break
		}

		key := pairs[i]
		value := pairs[i+1]

		// Register the required header
		r.headers[key] = value
	}

	return r
}

// CanRoute determines if the incoming request should use this route.
func (r *Route) canRoute(req events.APIGatewayProxyRequest) bool {
	// Loop through the expected headers & values
	for expKey, expValue := range r.headers {
		// If the header key is no present, we don't support this request
		if _, ok := req.Headers[expKey]; !ok {
			return false
		}

		// If the value is not what we expect from this key, we don't support
		// this request.
		if value := req.Headers[expKey]; value != expValue {
			return false
		}
	}

	// If the request method does not match this route's method, we don't
	// support this request.
	if req.HTTPMethod != r.method {
		return false
	}

	return true
}
