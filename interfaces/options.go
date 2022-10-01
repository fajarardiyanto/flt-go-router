package interfaces

import (
	loggerInterfaces "github.com/fajarardiyanto/flt-go-logger/interfaces"
	log "github.com/fajarardiyanto/flt-go-logger/lib"
	"net/http"
)

var (
	logger loggerInterfaces.Logger

	// DefaultPattern is pattern params query
	DefaultPattern = `[\w._-]+`

	METHOD = map[string]struct{}{
		http.MethodGet:    {},
		http.MethodPost:   {},
		http.MethodPut:    {},
		http.MethodDelete: {},
		http.MethodPatch:  {},
	}
)

func init() {
	logger = log.NewLib()
	logger.Init("HTTP Router")
	logger.SetOutputFormat(loggerInterfaces.OutputFormatDefault)
}

func GetLogger() loggerInterfaces.Logger {
	return logger
}
