package middleware

import (
	"github.com/alauda/bergamot/log"
	"github.com/alauda/loggo"
)

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

// LogFunc defined type
type LogFunc func(message string, fields loggo.Fields)

// GetStFunc get a logger structured function function from base logger
func (mw BaseLogger) GetStFunc(err error) LogFunc {
	return GetStFunc(err, mw.logger)

}

// GetStFunc open function to get a logger function based on the error
func GetStFunc(err error, logger log.Logger) LogFunc {
	if err != nil {
		return logger.StError
	}
	return logger.StInfo
}
