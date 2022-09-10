package interfaces

import (
	"net/http"
)

type MiddlewareFunc func(next Handler) Handler
type Handler func(rw http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_ = h(rw, r)
}
