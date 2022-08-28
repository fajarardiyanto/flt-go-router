package main

import (
	"github.com/fajarardiyanto/flt-go-router/interfaces"
	"github.com/fajarardiyanto/flt-go-router/lib"
	"net/http"
)

func main() {
	r := lib.New()
	r.Use(MiddlewareLogger(), MiddlewareError())

	r.GET("/ping", func(w http.ResponseWriter, r *http.Request) error {
		name := interfaces.GetQuery(r, "name")
		return interfaces.JSON(w, http.StatusOK, map[string]interface{}{
			"query": name,
		})
	})

	r.Run("8080")
}
