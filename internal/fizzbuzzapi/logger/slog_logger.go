package logger

import (
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	logger := slog.Default()

	// TODO: make log level configurable
	slog.SetLogLoggerLevel(slog.LevelInfo)
	return &SlogLogger{
		logger: logger,
	}
}

func (s *SlogLogger) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}
