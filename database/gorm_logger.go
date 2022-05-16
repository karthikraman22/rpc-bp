package database

import (
	"context"
	"time"

	"github.com/karthikraman22/rpc-bp/logger"
	"go.uber.org/zap"
	gorm_logger "gorm.io/gorm/logger"
)

type GormLogger struct {
	log logger.Logger
}

func NewGormLogger() *GormLogger {
	return &GormLogger{log: logger.WithNameOptions("gorm", zap.AddCaller(), zap.AddCallerSkip(2))}
}

func (gl *GormLogger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	newlogger := *gl
	return &newlogger
}

func (gl *GormLogger) Info(ctx context.Context, m string, v ...interface{}) {
	gl.log.Info(m, v)
}
func (gl *GormLogger) Warn(ctx context.Context, m string, v ...interface{}) {
	gl.log.Warn(m, v)
}

func (gl *GormLogger) Error(ctx context.Context, m string, v ...interface{}) {
	gl.log.Error(v[0].(error), m)
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		gl.log.Error(err, "trace", "sql", sql, "rows_affected", rows, "elapsed", time.Duration(elapsed))
	} else {
		gl.log.Info("trace", "sql", sql, "rows_affected", rows, "elapsed", time.Duration(elapsed))
	}
}
