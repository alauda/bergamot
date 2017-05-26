package middleware

import "github.com/alauda/bergamot/log"

// based on: https://gokit.io/examples/stringsvc.html

// BaseLogger middleware for logging
type BaseLogger struct {
	logger log.Logger
}

// NewLog constructor
func NewLog(logger log.Logger) BaseLogger {
	return BaseLogger{
		logger: logger,
	}
}
