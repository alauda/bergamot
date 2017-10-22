// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

// Fields for structured logging
type Fields map[string]interface{}

// A Logger represents a logging module. It has an associated logging
// level which can be changed; messages of lesser severity will
// be dropped. Loggers have a hierarchical relationship - see
// the package documentation.
//
// The zero Logger value is usable - any messages logged
// to it will be sent to the root Logger.
type Logger struct {
	impl *module
	size int
}

func (logger Logger) getModule() *module {
	if logger.impl == nil {
		return defaultContext.root
	}
	return logger.impl
}

// Name returns the logger's module name.
func (logger Logger) Name() string {
	return logger.getModule().Name()
}

// LogLevel returns the configured min log level of the logger.
func (logger Logger) LogLevel() Level {
	return logger.getModule().level
}

// EffectiveLogLevel returns the effective min log level of
// the receiver - that is, messages with a lesser severity
// level will be discarded.
//
// If the log level of the receiver is unspecified,
// it will be taken from the effective log level of its
// parent.
func (logger Logger) EffectiveLogLevel() Level {
	return logger.getModule().getEffectiveLogLevel()
}

// SetLogLevel sets the severity level of the given logger.
// The root logger cannot be set to UNSPECIFIED level.
// See EffectiveLogLevel for how this affects the
// actual messages logged.
func (logger Logger) SetLogLevel(level Level) {
	logger.getModule().setLevel(level)
}

// Logf logs a printf-formatted message at the given level.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger Logger) Logf(level Level, message string, args ...interface{}) {
	logger.LogCallf(2, level, message, args...)
}

// LogCallf logs a printf-formatted message at the given level.
// The location of the call is indicated by the calldepth argument.
// A calldepth of 1 means the function that called this function.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger Logger) LogCallf(calldepth int, level Level, message string, args ...interface{}) {
	module := logger.getModule()
	if !module.willWrite(level) {
		return
	}
	// Gather time, and filename, line number.
	now := time.Now() // get this early.
	// Param to Caller is the call depth.  Since this method is called from
	// the Logger methods, we want the place that those were called from.
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	}
	// Trim newline off format string, following usual
	// Go logging conventions.
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[0 : len(message)-1]
	}

	// To avoid having a proliferation of Info/Infof methods,
	// only use Sprintf if there are any args, and rely on the
	// `go vet` tool for the obvious cases where someone has forgotten
	// to provide an arg.
	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}
	module.write(Entry{
		Level:     level,
		Filename:  file,
		Line:      line,
		Timestamp: now,
		Message:   formattedMessage,
	})
}

// Criticalf logs the printf-formatted message at critical level.
func (logger Logger) Criticalf(message string, args ...interface{}) {
	logger.Logf(CRITICAL, message, args...)
}

// Errorf logs the printf-formatted message at error level.
func (logger Logger) Errorf(message string, args ...interface{}) {
	logger.Logf(ERROR, message, args...)
}

// Warningf logs the printf-formatted message at warning level.
func (logger Logger) Warningf(message string, args ...interface{}) {
	logger.Logf(WARNING, message, args...)
}

// Infof logs the printf-formatted message at info level.
func (logger Logger) Infof(message string, args ...interface{}) {
	logger.Logf(INFO, message, args...)
}

// Debugf logs the printf-formatted message at debug level.
func (logger Logger) Debugf(message string, args ...interface{}) {
	logger.Logf(DEBUG, message, args...)
}

// Tracef logs the printf-formatted message at trace level.
func (logger Logger) Tracef(message string, args ...interface{}) {
	logger.Logf(TRACE, message, args...)
}

func limitString(message string, size int, char rune) string {
	if len(message) == size {
		return message
	} else if len(message) > size {
		return message[:size]
	} else {
		runes := make([]byte, size)
		copy(runes, []byte(message))
		byteRune := byte(char)
		for i := len(message); i < size; i++ {
			runes[i] = byteRune
		}
		return string(runes)
	}
}

func (logger Logger) generateStructured(message string, fields Fields) (format string, args []interface{}) {
	var buf bytes.Buffer
	if len(message) > 0 {
		fields["msg"] = message
	}
	args = make([]interface{}, len(fields)*2)
	index := 0
	for k, v := range fields {
		buf.WriteString("%s")
		args[index] = k
		index++

		buf.WriteString("=")

		buf.WriteString("\"%v\" ")
		args[index] = v
		index++
	}
	format = buf.String()
	return
}

// StCritical structured logs at critical level.
func (logger Logger) StCritical(message string, fields Fields) {
	format, args := logger.generateStructured(message, fields)
	logger.Logf(CRITICAL, format, args...)
}

// StError structured logs at error level.
func (logger Logger) StError(message string, fields Fields) {
	format, args := logger.generateStructured(message, fields)
	logger.Logf(ERROR, format, args...)
}

// StWarning structured logs at  warning level.
func (logger Logger) StWarning(message string, fields Fields) {
	format, args := logger.generateStructured(message, fields)
	logger.Logf(WARNING, format, args...)
}

// StInfo structured logs at  info level.
func (logger Logger) StInfo(message string, fields Fields) {
	format, args := logger.generateStructured(message, fields)
	logger.Logf(INFO, format, args...)
}

// StDebug structured logs at  debug level.
func (logger Logger) StDebug(message string, fields Fields) {
	format, args := logger.generateStructured(message, fields)
	logger.Logf(DEBUG, format, args...)
}

// StTrace structured logs at trace level.
func (logger Logger) StTrace(message string, fields Fields) {
	format, args := logger.generateStructured(message, fields)
	logger.Logf(TRACE, format, args...)
}

// IsLevelEnabled returns whether debugging is enabled
// for the given log level.
func (logger Logger) IsLevelEnabled(level Level) bool {
	return logger.getModule().willWrite(level)
}

// IsErrorEnabled returns whether debugging is enabled
// at error level.
func (logger Logger) IsErrorEnabled() bool {
	return logger.IsLevelEnabled(ERROR)
}

// IsWarningEnabled returns whether debugging is enabled
// at warning level.
func (logger Logger) IsWarningEnabled() bool {
	return logger.IsLevelEnabled(WARNING)
}

// IsInfoEnabled returns whether debugging is enabled
// at info level.
func (logger Logger) IsInfoEnabled() bool {
	return logger.IsLevelEnabled(INFO)
}

// IsDebugEnabled returns whether debugging is enabled
// at debug level.
func (logger Logger) IsDebugEnabled() bool {
	return logger.IsLevelEnabled(DEBUG)
}

// IsTraceEnabled returns whether debugging is enabled
// at trace level.
func (logger Logger) IsTraceEnabled() bool {
	return logger.IsLevelEnabled(TRACE)
}
