package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/karthikraman22/rpc-bp/operation"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = logr.Logger

var (
	zapLog        *zap.Logger
	operationMode operation.OperationMode
	logMode       string
)

const (
	DebugLevel = 1
	InfoLevel  = 0
	WarnLevel  = -1
	ErrorLevel = -2
)

func init() {
	var err error

	// from build flag
	if logMode == "noop" {
		zapLog = zap.NewNop()
		return
	}

	// from env
	operationMode = operation.GetOperationMode()
	switch operationMode {
	case operation.DEVELOPMENT:
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.DisableStacktrace = true
		zapLog, err = cfg.Build()
	case operation.RELEASE:
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.DisableStacktrace = true
		zapLog, err = cfg.Build()
	default:
		zapLog = zap.NewNop()
	}

	if err != nil {
		panic(fmt.Sprintf("failed to setup logging: %v", err))
	}
}

// Creates a new logger with the given options
func WithOptions(opts ...zap.Option) logr.Logger {
	switch operationMode {
	case operation.DEVELOPMENT:
		return zapr.NewLogger(zapLog.WithOptions(opts...)).V(DebugLevel)
	case operation.RELEASE:
		return zapr.NewLogger(zapLog.WithOptions(opts...)).V(InfoLevel)
	default:
		return zapr.NewLogger(zapLog.WithOptions(opts...)).V(ErrorLevel)
	}
}

// Creates a new logger with the given name
func WithName(name string) logr.Logger {
	return WithOptions(zap.AddCaller()).WithName(name)
}

func UnderLying() *zap.Logger {
	return zapLog
}
