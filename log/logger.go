package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alauda/bergamot/contexts"
	"github.com/alauda/bergamot/loggo"
	aloggo "github.com/alauda/loggo"
)

// StandardLogger stan
type StandardLogger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// StructuredLogger simplified structured logger signature
type StructuredLogger interface {
	Trace(msg string, keyValues ...interface{})
	Debug(msg string, keyValues ...interface{})
	Info(msg string, keyValues ...interface{})
	Warning(msg string, keyValues ...interface{})
	Error(msg string, keyValues ...interface{})
}

// StLogger St logging functions for using with loggo.Fields
type StLogger interface {
	StCritical(message string, fields loggo.Fields)
	StError(message string, fields loggo.Fields)
	StWarning(message string, fields loggo.Fields)
	StInfo(message string, fields loggo.Fields)
	StDebug(message string, fields loggo.Fields)
	StTrace(message string, fields loggo.Fields)
}

// Logger interface to define a logger entity
type Logger interface {
	StandardLogger
	StructuredLogger
	StLogger
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
// @param: depth given one integer will set the depth in which will find during runtime the file and line number
// recommended to use 3 as a number (only the first integer will be used)
func NewLogger(packageName string, size ...int) Logger {
	return loggo.GetLogger(packageName, size...)
}

// NewStandardLogger constructs a standard logger
func NewStandardLogger(packageName string, depth ...int) StandardLogger {
	return aloggo.GetLogger(packageName, depth...)
}

// NewNoCodeLogger will not print file
func NewNoCodeLogger(name string, size ...int) Logger {
	writter := loggo.NewSimpleWriter(os.Stdout, NoCodeFormat)
	loggo.ReplaceDefaultWriter(writter)
	return loggo.GetLogger(name, size...)
}

// NoCodeFormat will remove the code reference from msg
func NoCodeFormat(entry loggo.Entry) string {
	ts := entry.Timestamp.In(time.UTC).Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s %s %s %s", ts, entry.Level, entry.Module, entry.Message)
}

// F generates a map using a key-value pairs as arguments
// for example, if given arguments are "key" and 1 will return
// map[string]interface{}{"key": 1}
// PS: will not add key if the value is missing
func F(args ...interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	if len(args) > 0 {
		// avoiding nil cases
		var (
			k  string
			ok bool
		)
		for i, v := range args {
			if i%2 == 0 {
				if k, ok = v.(string); !ok {
					k = fmt.Sprintf("%v", v)
				}
			} else {
				res[k] = v
			}
		}
	}
	return
}
