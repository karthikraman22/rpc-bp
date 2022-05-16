package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	name   string
	logger *zap.Logger
}

func newZapLoggerWithOptions(name string, options ...zap.Option) *zapLogger {

	logCfg := zap.NewProductionConfig()
	logCfg.DisableStacktrace = true
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := logCfg.Build()

	if err != nil {
		panic(fmt.Sprintf("failed to setup logging: %v", err))
	}

	return &zapLogger{name: name, logger: logger.WithOptions(options...).Named(name)}
}

func newZapLogger(name string) *zapLogger {
	return newZapLoggerWithOptions(name, zap.AddCaller(), zap.AddCallerSkip(1))
}

// Info logs a message at level Info.
func (zl *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	if checkedEntry := zl.logger.Check(zap.InfoLevel, msg); checkedEntry != nil {
		checkedEntry.Write(zl.handleFields(keysAndValues)...)
	}
}

// Debug logs a message at level Debug.
func (zl *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	if checkedEntry := zl.logger.Check(zap.DebugLevel, msg); checkedEntry != nil {
		checkedEntry.Write(zl.handleFields(keysAndValues)...)
	}
}

// Warn logs a message at level Warn.
func (zl *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	if checkedEntry := zl.logger.Check(zap.WarnLevel, msg); checkedEntry != nil {
		checkedEntry.Write(zl.handleFields(keysAndValues)...)
	}
}

// Error logs a message at level Error.
func (zl *zapLogger) Error(errVal error, msg string, keysAndValues ...interface{}) {
	if checkedEntry := zl.logger.Check(zap.ErrorLevel, msg); checkedEntry != nil {
		checkedEntry.Write(zl.handleFields(keysAndValues, zap.NamedError("error", errVal))...)
	}
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (zl *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	if checkedEntry := zl.logger.Check(zap.FatalLevel, msg); checkedEntry != nil {
		checkedEntry.Write(zl.handleFields(keysAndValues)...)
	}
}

func (zl *zapLogger) handleFields(args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		// Slightly slower fast path when we need to inject "v".
		return additional
	}
	numFields := len(args)/2 + len(additional)
	fields := make([]zap.Field, 0, numFields)
	for i := 0; i < len(args); {
		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			zl.logger.WithOptions(zap.AddCallerSkip(1)).DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}
		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			zl.logger.WithOptions(zap.AddCallerSkip(1)).DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}
	return append(fields, additional...)
}
