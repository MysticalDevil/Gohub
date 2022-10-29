// Package logger Handle log related logic
package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gohub/app"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
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
