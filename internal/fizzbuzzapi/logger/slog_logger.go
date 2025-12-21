package logger

import (
	"log/slog"
)

func NewSlogLogger() *slog.Logger {
	logger := slog.Default()

	// TODO: make log level configurable
	slog.SetLogLoggerLevel(slog.LevelInfo)
	return logger
}
