package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init() {
	var cfg zap.Config
	env := os.Getenv("ENV")

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	logFormat := strings.ToLower(os.Getenv("LOG_FORMAT"))

	if env == "production" || logFormat == "json" {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
	}

	switch logLevel {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		if env == "production" {
			cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		} else {
			cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		}
	}

	cfg.OutputPaths = []string{"stdout"}      // logs INFO, DEBUG, etc -> (terminal)
	cfg.ErrorOutputPaths = []string{"stderr"} // logs ERROR, FATAL, etc -> (terminal error channel)
	Logger, _ = cfg.Build()
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Wrapper functions for the Logger

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
