package log

import "github.com/alauda/bergamot/loggo"

// EmptyLogger empty logger (does nothing)
type EmptyLogger struct{}

// Tracef nothing
func (EmptyLogger) Tracef(format string, args ...interface{}) {}

// Debugf nothing
func (EmptyLogger) Debugf(format string, args ...interface{}) {}

// Infof nothing
func (EmptyLogger) Infof(format string, args ...interface{}) {}

// Warningf nothing
func (EmptyLogger) Warningf(format string, args ...interface{}) {}

// Errorf nothing
func (EmptyLogger) Errorf(format string, args ...interface{}) {}

// Trace nothing
func (EmptyLogger) Trace(format string, args ...interface{}) {}

// Debug nothing
func (EmptyLogger) Debug(format string, args ...interface{}) {}

// Info nothing
func (EmptyLogger) Info(format string, args ...interface{}) {}

// Warning nothing
func (EmptyLogger) Warning(format string, args ...interface{}) {}

// Error nothing
func (EmptyLogger) Error(format string, args ...interface{}) {}

// StCritical nothing
func (EmptyLogger) StCritical(message string, fields loggo.Fields) {}

// StError nothing
func (EmptyLogger) StError(message string, fields loggo.Fields) {}

// StWarning nothing
func (EmptyLogger) StWarning(message string, fields loggo.Fields) {}

// StInfo nothing
func (EmptyLogger) StInfo(message string, fields loggo.Fields) {}

// StDebug nothing
func (EmptyLogger) StDebug(message string, fields loggo.Fields) {}

// StTrace nothing
func (EmptyLogger) StTrace(message string, fields loggo.Fields) {}

var emptyLogger = EmptyLogger{}

// GetSafe will verify if given logger is initiated and will return
// an empty logger when it is not to make log calls safe
func GetSafe(logger Logger) Logger {
	if logger == nil {
		return emptyLogger
	}
	return logger
}
