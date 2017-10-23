package log

import (
	"context"

	"github.com/alauda/bergamot/contexts"
	"github.com/alauda/loggo"
)

// Logger interface to define a logger entity
type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	StCritical(message string, fields loggo.Fields)
	StError(message string, fields loggo.Fields)
	StWarning(message string, fields loggo.Fields)
	StInfo(message string, fields loggo.Fields)
	StDebug(message string, fields loggo.Fields)
	StTrace(message string, fields loggo.Fields)
}

// Level defines log levels
type Level int

const (
	// LevelTrace prints all tracing and up
	LevelTrace Level = iota
	// LevelDebug prints all messages
	LevelDebug
	// LevelInfo prints only info or higher
	LevelInfo
	// LevelError prints only error messages
	LevelError
)

// SetLevel Sets logging levels
func SetLevel(level Level) {
	var config string
	switch level {

	case LevelInfo:
		config = "<root>=INFO"
	case LevelError:
		config = "<root>=ERROR"
	case LevelDebug:
		config = "<root>=DEBUG"
	case LevelTrace:
		fallthrough
	default:
		config = "<root>=TRACE"

	}
	loggo.ConfigureLoggers(config)
}

// GetFields get fields using a context
func GetFields(ctx context.Context) (fields loggo.Fields) {
	fields = loggo.Fields{}
	return AddRequestID(ctx, fields)
}

// AddRequestID adds the request Id to fields if one is there
func AddRequestID(ctx context.Context, fields loggo.Fields) loggo.Fields {
	if ctx != nil {
		if request := contexts.GetRequestID(ctx); request != "" {
			fields["request_id"] = request
		}
	}
	return fields
}

// NewLogger constructs a new logger for package
// @param: packageName stands for the package of your logger
// @param: size given one integer will set the size for the message on structured log
func NewLogger(packageName string, size ...int) Logger {
	return loggo.GetLogger(packageName, size...)
}
