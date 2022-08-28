package main

import (
	loggerInterfaces "github.com/fajarardiyanto/flt-go-logger/interfaces"
	log "github.com/fajarardiyanto/flt-go-logger/lib"
)

var logger loggerInterfaces.Logger

func init() {
	logger = log.NewLib()
	logger.Init("HTTP Router")
}

func GetLogger() loggerInterfaces.Logger {
	return logger
}
