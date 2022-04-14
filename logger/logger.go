package logger

import (
	"fmt"

	"achuala.in/rpc-bp/operation"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
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
		zapLog, err = zap.NewDevelopment()
	case operation.RELEASE:
		zapLog, err = zap.NewProduction()
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
