// Package logger Handle log related logic
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gohub/pkg/app"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger Global logger object
var Logger *zap.Logger

// InitLogger Log initialization
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {
	// Get log write medium
	writerSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)

	// Set log level
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("Log initialization error, log level setting is wrong. " +
			"Please modify the log.level configuration item in the config/log.go file")
	}

	// Initialize core
	core := zapcore.NewCore(getEncoder(), writerSyncer, logLevel)

	// Initialize Logger
	Logger = zap.New(core,
		zap.AddCaller(),                   // Call file and line number, internally use runtime.Caller
		zap.AddCallerSkip(1),              // One layer is encapsulated, and the calling file removes one layer (runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // The stacktrace will only be displayed when there is an Error
	)

	// Replace custom logger with global logger
	// When zap.L().Fatal() is called, our custom Logger will be used
	zap.ReplaceGlobals(Logger)
}

// getEncoder Set the log storage format
func getEncoder() zapcore.Encoder {
	// log format rules
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,        // Add '\n' to each log line
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Log level name uppercase
		EncodeTime:     customTimeEncoder,                // Time format, 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder,   // Execution time, in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,       // Caller short format, such as types/converter.go:17
	}

	// local environment configuration
	if app.IsLocal() {
		// Keyword highlighting in terminal output
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// Locally set the built-in Console decoder (support stacktrace wrapping)
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Online environment using JSON decoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// customTimeEncoder Custom friendly time format
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter Logging medium
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// If log file by date is configured
	if logType == "daily" {
		logName := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logName)
	}

	// rolling log
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	// Configure output media
	if app.IsLocal() {
		// Terminal prints and logs files when developing locally
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		// Only log files in production environment
		return zapcore.AddSync(lumberJackLogger)
	}
}

// Dump is dedicated to debugging, will not interrupt the program, and will print out warning messages on the terminal
// The first parameter will use json.Marshal for rendering, the second parameter message (optional)
//
//	logger.Dump(user.User{Name"test"})
//	logger.Dump(user.User(Name:"test"}, "User info")
func Dump(value any, msg ...string) {
	valueString := jsonString(value)
	if len(msg) > 0 {
		Logger.Warn("Dump", zap.String(msg[0], valueString))
	} else {
		Logger.Warn("Dump", zap.String("data", valueString))
	}
}

// LogIf When err != nil log error level log
func LogIf(err error) {
	if err != nil {
		Logger.Error("Error Occurred:", zap.Error(err))
	}
}

// LogWarnIf When err != nil log warning level log
func LogWarnIf(err error) {
	if err != nil {
		Logger.Warn("Error Occurred:", zap.Error(err))
	}
}

// LogInfoIf When err != nil log info level log
func LogInfoIf(err error) {
	if err != nil {
		Logger.Info("Error Occurred:", zap.Error(err))
	}
}

// Debug log, detailed program log
// e.g.:
//
//	logger.Debug("Database", zap.String("sql", sql))
func Debug(moduleName string, fields ...zap.Field) {
	Logger.Debug(moduleName, fields...)
}

// Info log, notification log
func Info(moduleName string, fields ...zap.Field) {
	Logger.Info(moduleName, fields...)
}

// Warn log, warning log
func Warn(moduleName string, fields ...zap.Field) {
	Logger.Warn(moduleName, fields...)
}

// Error log, log on error, should not interrupt the program
func Error(moduleName string, fields ...zap.Field) {
	Logger.Error(moduleName, fields...)
}

// Fatal log, after writing the log, call os.Exit(1) to exit the program
func Fatal(moduleName string, fields ...zap.Field) {
	Logger.Fatal(moduleName, fields...)
}

// DebugString Record a debug log of string type
// e.g.
//
//	logger.DebugString("SMS", "SMS content", string(result.RawResponse))
func DebugString(moduleName, name, msg string) {
	Logger.Debug(moduleName, zap.String(name, msg))
}

func InfoString(moduleName, name, msg string) {
	Logger.Info(moduleName, zap.String(name, msg))
}

func WarnString(moduleName, name, msg string) {
	Logger.Warn(moduleName, zap.String(name, msg))
}

func ErrorString(moduleName, name, msg string) {
	Logger.Error(moduleName, zap.String(name, msg))
}

func FatalString(moduleName, name, msg string) {
	Logger.Fatal(moduleName, zap.String(name, msg))
}

// DebugJSON Record debug logs for object types
// e.g.
//
//	logger.DebugJSON("Auth", "Read log in users", auth.CurrentUser())
func DebugJSON(moduleName, name string, value any) {
	Logger.Debug(moduleName, zap.String(name, jsonString(value)))
}

func InfoJSON(moduleName, name string, value any) {
	Logger.Info(moduleName, zap.String(name, jsonString(value)))
}

func WarnJSON(moduleName, name string, value any) {
	Logger.Warn(moduleName, zap.String(name, jsonString(value)))
}

func ErrorJSON(moduleName, name string, value any) {
	Logger.Error(moduleName, zap.String(name, jsonString(value)))
}

func FatalJSON(moduleName, name string, value any) {
	Logger.Fatal(moduleName, zap.String(name, jsonString(value)))
}

func jsonString(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		Logger.Error("Logger", zap.String("JSON marshal error", err.Error()))
	}
	return string(b)
}
