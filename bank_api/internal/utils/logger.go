package utils

import (
	"log/slog"
	"os"
)

// InitializeLogger sets up the logger configuration.
func InitializeLogger(config *Config) {
	leveler := new(slog.LevelVar)
	leveler.Set(translateLogLevel(config.LogLevel))
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: leveler})
	slog.SetDefault(slog.New(handler))
}

func translateLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func LogFatal(msg string, err error) {
	slog.Error(msg, slog.Any("error", err))
	os.Exit(1)
}
