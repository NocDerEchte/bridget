package logging

import (
	"log/slog"
	"os"

	"github.com/nocderechte/bridget/pkg/helper"
)

type Logger struct {
	logFormat string
	logLevel  string
	logger    *slog.Logger
}

func NewLogger() *Logger {
	l := Logger{}

	l.logFormat = helper.GetEnv("LOG_FORMAT", "json") // json (default), text
	l.logLevel = helper.GetEnv("LOG_LEVEL", "info")   // debug, info (default), warn, error

	opts := &slog.HandlerOptions{}

	switch l.logLevel {
	case "debug":
		opts.Level = slog.LevelDebug
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
		l.logLevel = "info"
	}

	var handler slog.Handler

	if l.logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l.logger = slog.New(handler)
	l.logger.Debug("Successfully initialized logging system. Level: " + l.logLevel + ", format: " + l.logFormat)

	return &l
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Errorf(msg string, err error, args ...any) {
	l.logger.With(slog.String("error", err.Error())).Error(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
