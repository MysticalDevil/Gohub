// Package logger Handle log related logic
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"gohub/pkg/app"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger Global logger object
var Logger *slog.Logger

// InitLogger Log initialization
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {
	writer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)
	logLevel, ok := parseLevel(level)
	if !ok {
		fmt.Println("Log initialization error, log level setting is wrong. " +
			"Please modify the log.level configuration item in the config/log.go file")
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	var handler slog.Handler
	if app.IsLocal() {
		handler = slog.NewTextHandler(writer, handlerOpts)
	} else {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

func parseLevel(level string) (slog.Level, bool) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

// getLogWriter Logging medium
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) io.Writer {
	if logType == "daily" {
		logName := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logName)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	if app.IsLocal() {
		return io.MultiWriter(os.Stdout, lumberJackLogger)
	}
	return lumberJackLogger
}

// Dump is dedicated to debugging, will not interrupt the program, and will print out warning messages on the terminal
func Dump(value any, msg ...string) {
	valueString := jsonString(value)
	if len(msg) > 0 {
		Logger.Warn("Dump", slog.String(msg[0], valueString))
	} else {
		Logger.Warn("Dump", slog.String("data", valueString))
	}
}

// LogIf When err != nil log error level log
func LogIf(err error) {
	if err != nil {
		Logger.Error("Error Occurred", slog.Any("error", err))
	}
}

// LogWarnIf When err != nil log warning level log
func LogWarnIf(err error) {
	if err != nil {
		Logger.Warn("Error Occurred", slog.Any("error", err))
	}
}

// LogInfoIf When err != nil log info level log
func LogInfoIf(err error) {
	if err != nil {
		Logger.Info("Error Occurred", slog.Any("error", err))
	}
}

// Debug log, detailed program log
func Debug(moduleName string, fields ...slog.Attr) {
	Logger.Debug(moduleName, attrsToAny(fields)...)
}

// Info log, notification log
func Info(moduleName string, fields ...slog.Attr) {
	Logger.Info(moduleName, attrsToAny(fields)...)
}

// Warn log, warning log
func Warn(moduleName string, fields ...slog.Attr) {
	Logger.Warn(moduleName, attrsToAny(fields)...)
}

// Error log, log on error, should not interrupt the program
func Error(moduleName string, fields ...slog.Attr) {
	Logger.Error(moduleName, attrsToAny(fields)...)
}

// Fatal log, after writing the log, call os.Exit(1) to exit the program
func Fatal(moduleName string, fields ...slog.Attr) {
	Logger.Error(moduleName, attrsToAny(fields)...)
	os.Exit(1)
}

// DebugString Record a debug log of string type
func DebugString(moduleName, name, msg string) {
	Logger.Debug(moduleName, slog.String(name, msg))
}

func InfoString(moduleName, name, msg string) {
	Logger.Info(moduleName, slog.String(name, msg))
}

func WarnString(moduleName, name, msg string) {
	Logger.Warn(moduleName, slog.String(name, msg))
}

func ErrorString(moduleName, name, msg string) {
	Logger.Error(moduleName, slog.String(name, msg))
}

func FatalString(moduleName, name, msg string) {
	Logger.Error(moduleName, slog.String(name, msg))
	os.Exit(1)
}

// DebugJSON Record debug logs for object types
func DebugJSON(moduleName, name string, value any) {
	Logger.Debug(moduleName, slog.String(name, jsonString(value)))
}

func InfoJSON(moduleName, name string, value any) {
	Logger.Info(moduleName, slog.String(name, jsonString(value)))
}

func WarnJSON(moduleName, name string, value any) {
	Logger.Warn(moduleName, slog.String(name, jsonString(value)))
}

func ErrorJSON(moduleName, name string, value any) {
	Logger.Error(moduleName, slog.String(name, jsonString(value)))
}

func FatalJSON(moduleName, name string, value any) {
	Logger.Error(moduleName, slog.String(name, jsonString(value)))
	os.Exit(1)
}

func jsonString(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		Logger.Error("Logger", slog.String("JSON marshal error", err.Error()))
	}
	return string(b)
}

func attrsToAny(attrs []slog.Attr) []any {
	if len(attrs) == 0 {
		return nil
	}
	out := make([]any, len(attrs))
	for i, attr := range attrs {
		out[i] = attr
	}
	return out
}
