package logger

import (
	"strings"
	"sync"
)

const (
	// LogTypeLog is normal log type.
	LogTypeLog = "log"
	// LogTypeRequest is Request log type.
	LogTypeRequest = "request"

	// Field names that defines Dapr log schema.
	logFieldTimeStamp = "time"
	logFieldLevel     = "level"
	logFieldType      = "type"
	logFieldScope     = "scope"
	logFieldMessage   = "msg"
	logFieldInstance  = "instance"
	logFieldDaprVer   = "ver"
	logFieldAppID     = "app_id"
)

// LogLevel is Dapr Logger Level type.
type LogLevel string

const (
	// DebugLevel has verbose message.
	DebugLevel LogLevel = "debug"
	// InfoLevel is default log level.
	InfoLevel LogLevel = "info"
	// WarnLevel is for logging messages about possible issues.
	WarnLevel LogLevel = "warn"
	// ErrorLevel is for logging errors.
	ErrorLevel LogLevel = "error"
	// FatalLevel is for logging fatal messages. The system shuts down after logging the message.
	FatalLevel LogLevel = "fatal"

	// UndefinedLevel is for undefined log level.
	UndefinedLevel LogLevel = "undefined"
)

// globalLoggers is the collection of Logger that is shared globally.
// TODO: User will disable or enable logger on demand.
var (
	globalLoggers     = map[string]Logger{}
	globalLoggersLock = sync.RWMutex{}
)

// Logger includes the logging api sets.
type Logger interface {
	// Info logs a message at level Info.
	Info(msg string, keysAndValues ...interface{})
	// Debug logs a message at level Debug.
	Debug(msg string, keysAndValues ...interface{})
	// Warn logs a message at level Warn.
	Warn(msg string, keysAndValues ...interface{})
	// Error logs a message at level Error.
	Error(errVal error, keysAndValues string, args ...interface{})
	// Fatal logs a message at level Fatal then the process will exit with status set to 1.
	Fatal(msg string, keysAndValues ...interface{})
}

// toLogLevel converts to LogLevel.
func toLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	}

	// unsupported log level by Dapr
	return UndefinedLevel
}

// NewLogger creates new Logger instance.
func WithName(name string) Logger {
	globalLoggersLock.Lock()
	defer globalLoggersLock.Unlock()

	logger, ok := globalLoggers[name]
	if !ok {
		logger = newZapLogger(name)
		globalLoggers[name] = logger
	}

	return logger
}

func getLoggers() map[string]Logger {
	globalLoggersLock.RLock()
	defer globalLoggersLock.RUnlock()

	l := map[string]Logger{}
	for k, v := range globalLoggers {
		l[k] = v
	}

	return l
}
