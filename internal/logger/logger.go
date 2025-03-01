package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/guarilha/go-ddd-starter/internal/config"
)

func NewLogger(cfg config.Config) *slog.Logger {
	logOutput := os.Stdout
	if cfg.Logging.Stderr {
		logOutput = os.Stderr
	}

	logOptions := slog.HandlerOptions{
		Level: getLogLevel(cfg),
	}

	logHandler := getLogHandler(cfg, logOutput, &logOptions)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return logger
}

func getLogLevel(cfg config.Config) slog.Level {
	// If a specific level is configured, use it
	if cfg.Logging.Level != "" {
		switch strings.ToUpper(cfg.Logging.Level) {
		case "DEBUG":
			return slog.LevelDebug
		case "INFO":
			return slog.LevelInfo
		case "WARN", "WARNING":
			return slog.LevelWarn
		case "ERROR":
			return slog.LevelError
		}
	}

	// Otherwise, use the level based on environment
	if cfg.Environment == "development" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func getLogHandler(cfg config.Config, output *os.File, opts *slog.HandlerOptions) slog.Handler {
	// If a specific type is configured, use it
	if cfg.Logging.Type != "" {
		switch strings.ToUpper(cfg.Logging.Type) {
		case "JSON":
			return slog.NewJSONHandler(output, opts)
		case "TEXT":
			return slog.NewTextHandler(output, opts)
		}
	}

	// Otherwise, use the type based on environment
	if cfg.Environment == "development" {
		return slog.NewTextHandler(output, opts)
	}
	return slog.NewJSONHandler(output, opts)
}
