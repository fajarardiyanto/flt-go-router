package main

import (
	"github.com/fajarardiyanto/flt-go-router/interfaces"
	"github.com/fajarardiyanto/flt-go-router/lib"
	"net/http"
)

func main() {
	r := lib.New("v1.0.0")
	r.Use(MiddlewareLogger(), MiddlewareError())

	g := r.Group("/group")
	g.GET("/ping", func(w http.ResponseWriter, r *http.Request) error {
		name := interfaces.GetQuery(r, "name")
		return interfaces.JSON(w, http.StatusOK, map[string]interface{}{
			"query": name,
		})
	})

	r.GET("/test", func(w http.ResponseWriter, r *http.Request) error {
		name := interfaces.GetQuery(r, "name")
		return interfaces.JSON(w, http.StatusOK, map[string]interface{}{
			"test": name,
		})
	})

	r.Run("8081")
}
