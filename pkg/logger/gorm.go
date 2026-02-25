package logger

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gohub/pkg/helpers"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// GormLogger operation object, implement gormLogger.Interface
type GormLogger struct {
	Logger        *slog.Logger
	SlowThreshold time.Duration
}

func NewGormLogger() GormLogger {
	return GormLogger{
		Logger:        Logger,
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LogMode Implement the LogMode method of gormLogger.Interface
func (l GormLogger) LogMode(_ gormLogger.LogLevel) gormLogger.Interface {
	return GormLogger{
		Logger:        l.Logger,
		SlowThreshold: l.SlowThreshold,
	}
}

// Info Implement the Info method of gormLogger.Interface
func (l GormLogger) Info(_ context.Context, str string, args ...any) {
	l.logger().Debug(str, slog.Any("args", args))
}

// Warn Implement the Warn method of gormLogger.Interface
func (l GormLogger) Warn(_ context.Context, str string, args ...any) {
	l.logger().Warn(str, slog.Any("args", args))
}

// Error Implement the Error method of gormLogger.Interface
func (l GormLogger) Error(_ context.Context, str string, args ...any) {
	l.logger().Error(str, slog.Any("args", args))
}

// Trace Implement the Trace method of gormLogger.Interface
func (l GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	logFields := []slog.Attr{
		slog.String("sql", sql),
		slog.String("time", helpers.MicrosecondStr(elapsed)),
		slog.Int64("rows", rows),
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrRecordNotFound", attrsToAny(logFields)...)
		} else {
			logFields = append(logFields, slog.Any("error", err))
			l.logger().Error("Database Error", attrsToAny(logFields)...)
		}
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger().Warn("Database Slow Log", attrsToAny(logFields)...)
	}

	l.logger().Debug("Database Query", attrsToAny(logFields)...)
}

func (l GormLogger) logger() *slog.Logger {
	if l.Logger != nil {
		return l.Logger
	}
	return slog.Default()
}
