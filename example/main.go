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
		file, err := interfaces.GetFormFile(r, "file")
		if err != nil {
			GetLogger().Error(err)
			return err
		}

		GetLogger().Info(file.Filename)
		return interfaces.JSON(w, http.StatusOK, map[string]interface{}{
			"query": name,
		})
	})

	r.GET("/test/:id/:guid", func(w http.ResponseWriter, r *http.Request) error {
		name := interfaces.GetFormData(r, "name")
		id := interfaces.GetParam(r, "id")
		guid := interfaces.GetParam(r, "guid")

		return interfaces.JSON(w, http.StatusOK, map[string]interface{}{
			"id":   id,
			"test": name,
			"guid": guid,
		})
	})

	r.Run("8081")
}
