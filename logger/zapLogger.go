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

func newZapLogger(name string) *zapLogger {

	logCfg := zap.NewProductionConfig()
	logCfg.DisableStacktrace = true
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := logCfg.Build(zap.AddCaller())

	if err != nil {
		panic(fmt.Sprintf("failed to setup logging: %v", err))
	}

	return &zapLogger{name: name, logger: logger.Named(name)}
}

// Info logs a message at level Info.
func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

// Infof logs a message at level Info.
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

// Debug logs a message at level Debug.
func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

// Debugf logs a message at level Debug.
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

// Warn logs a message at level Warn.
func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

// Warnf logs a message at level Warn.
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

// Error logs a message at level Error.
func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

// Errorf logs a message at level Error.
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}
