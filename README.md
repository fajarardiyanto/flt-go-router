### Go Module Router
```
  _____.__   __                       .___    .__          
_/ ____\  |_/  |_    _____   ____   __| _/_ __|  |   ____  
\   __\|  |\   __\  /     \ /  _ \ / __ |  |  \  | _/ __ \ 
 |  |  |  |_|  |   |  Y Y  (  <_> ) /_/ |  |  /  |_\  ___/ 
 |__|  |____/__|   |__|_|  /\____/\____ |____/|____/\___  >
                         \/            \/               \/   
```
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
go get -u github.com/fajarardiyanto/flt-go-router@v0.0.2
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
	r := lib.New("v1.0.0")
	r.Use(MiddlewareLogger(), MiddlewareError())

	g := r.Group("/group")
	g.GET("/ping", func(w http.ResponseWriter, r *http.Request) error {
		name := interfaces.GetQuery(r, "name")
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