package logging

import (
	"log/slog"
	"os"

	"github.com/nocderechte/bridget/pkg/helper"
)

var (
	logFormat = helper.GetEnv("LOG_FORMAT", "json") // json (default), text, dev
	logLevel  = helper.GetEnv("LOG_LEVEL", "info")  // debug, info (default), warn, error
	logger    *slog.Logger
)

func InitLogger() {
	opts := &slog.HandlerOptions{}

	switch logLevel {
	case "debug":
		opts.Level = slog.LevelDebug
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
		logLevel = "info"
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "dev":
		handler = &prettyHandler{level: opts.Level, out: os.Stdout}
	}

	logger = slog.New(handler)
	logger.Debug("Successfully initialized logging system. Level: " + logLevel + ", format: " + logFormat)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Errorf(msg string, err error, args ...any) {
	logger.With(slog.String("error", err.Error())).Error(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
