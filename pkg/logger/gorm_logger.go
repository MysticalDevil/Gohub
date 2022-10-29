package logger

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gohub/pkg/helpers"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GormLogger operation object, implement gormLogger.Interface
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     Logger,                 // Use the global logger.Logger object
		SlowThreshold: 200 * time.Millisecond, // Slow query threshold, in thousandths of a second
	}
}

// LogMode Implement the LogMode method of gormLogger.Interface
func (l GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return GormLogger{
		ZapLogger:     l.ZapLogger,
		SlowThreshold: l.SlowThreshold,
	}
}

// Info Implement the Info method of gormLogger.Interface
func (l GormLogger) Info(_ context.Context, str string, args ...any) {
	l.logger().Sugar().Debugf(str, args...)
}

// Warn Implement the Warn method of gormLogger.Interface
func (l GormLogger) Warn(_ context.Context, str string, args ...any) {
	l.logger().Sugar().Warnf(str, args...)
}

// Error Implement the Error method of gormLogger.Interface
func (l GormLogger) Error(_ context.Context, str string, args ...any) {
	l.logger().Sugar().Errorf(str, args...)
}

// Trace Implement the Trace method of gormLogger.Interface
func (l GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Get running time
	elapsed := time.Since(begin)
	// Get the number of SQL requests and returned items
	sql, rows := fc()

	// Common field
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", helpers.MicrosecondStr(elapsed)),
		zap.Int64("rows", rows),
	}

	// Gorm error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, zap.Error(err))
			l.logger().Error("Database Error", logFields...)
		}
	}

	// Slow query log
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger().Warn("Database Slow Log", logFields...)
	}

	// Log all SQL requests
	l.logger().Debug("Database Query", logFields...)
}

// logger Internal auxiliary method to ensure the accuracy of Zap's built-in information Caller
func (l GormLogger) logger() *zap.Logger {
	// skip calls to gorm builtins
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapGormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	// Subtract one encapsulation, and add zap.AddCallerSkip(1) to logger initialization once
	clone := l.ZapLogger.WithOptions(zap.AddCallerSkip(-2))

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapGormPackage):
		default:
			return clone.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}
