### Go Module Router

### Installation
```sh
go get github.com/fajarardiyanto/flt-go-router
```

###### Upgrading to the latest version
```sh
go get -u github.com/fajarardiyanto/flt-go-router
```

###### Upgrade or downgrade with tag version if available
```sh
go get -u github.com/fajarardiyanto/flt-go-router@v0.0.1
```

### Usage
```go
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
```

#### Run Example
```sh
make run
```

#### Tips
Maybe it would be better to do some basic code scanning before pushing to the repository.
```sh
# for *.nix users just run gosec.sh
# curl is required
# more information https://github.com/securego/gosec
make scan
```