package log

import (
	"log/slog"
	"os"
)

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func Info(msg string, keyValues ...any) {
	defaultLogger.Info(msg, keyValues...)
}

func Error(err error, msg string, keyValues ...any) {
	defaultLogger.With("error", err).Error(msg, keyValues...)
}

func Warn(msg string, keyValues ...any) {
	defaultLogger.Warn(msg, keyValues...)
}

func SetDefault(l *slog.Logger) {
	defaultLogger = l
}
