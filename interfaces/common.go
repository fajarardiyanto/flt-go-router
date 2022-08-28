package interfaces

import (
	loggerInterfaces "github.com/fajarardiyanto/flt-go-logger/interfaces"
	"net/http"
	"os"
	"strings"
)

type (
	// ContextKeyType is a private struct that is used for storing values in net.Context
	ContextKeyType struct{}

	// ParamsMapType is a private type that is used to store route params
	ParamsMapType map[string]string
)

var (
	// ContextKey is the key that is used to store values in net.Context for each request
	ContextKey = ContextKeyType{}
)

// GetAllParams returns all route params stored in http.Request.
func GetAllParams(r *http.Request) ParamsMapType {
	if values, ok := r.Context().Value(ContextKey).(ParamsMapType); ok {
		return values
	}
	return nil
}

func ResolveAddress(addr []string, logger loggerInterfaces.Logger) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logger.Info("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		logger.Info("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		var port string
		if strings.Contains(addr[0], ":") {
			port = addr[0]
		} else {
			port = ":" + addr[0]
		}

		return port
	default:
		panic("too many parameters")
	}
}
