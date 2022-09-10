package interfaces

import (
	"net/http"
)

type Routers interface {
	// GET adds the route path
	GET(path string, handle Handler)

	// POST adds the route path
	POST(path string, handle Handler)

	// DELETE adds the route path
	DELETE(path string, handle Handler)

	// PUT adds the route path
	PUT(path string, handle Handler)

	// PATCH adds the route path
	PATCH(path string, handle Handler)

	// Group define routes groups
	Group(prefix string) Routers

	//Run attaches the router to a http.Server and starts listening and serving HTTP requests.
	Run(addr ...string)

	// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
	RunTLS(addr, certFile, keyFile string) (err error)

	// Handle register a new request handler with the given path and method.
	Handle(method string, path string, handle Handler)

	// NotFoundFunc registers a handler when the request route is not found
	NotFoundFunc(handler Handler)

	// Use appends a Middleware handler to the Middleware stack.
	Use(middleware ...MiddlewareFunc)

	// HandleNotFound registers a handler when the request route is not found
	HandleNotFound(w http.ResponseWriter, req *http.Request, msg string, code int, middleware []MiddlewareFunc)

	// Match checks if the request matches the route pattern
	//Match(requestUrl string, path string) bool

	// MatchAndParse checks if the request matches the route path and returns a map of the parsed
	MatchAndParse(requestUrl string, path string) (matchParams ParamsMapType, b bool)

	// ServeHTTP makes the router implement the http.Handler interface.
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
