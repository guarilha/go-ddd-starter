package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/guarilha/go-ddd-starter/internal/config"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		config      config.Config
		checkLevel  slog.Level
		checkType   string
		checkOutput *os.File
	}{
		{
			name: "Development environment with default config",
			config: config.Config{
				Environment: "development",
				Logging:     config.LoggingConfig{},
			},
			checkLevel:  slog.LevelDebug,
			checkType:   "text",
			checkOutput: os.Stdout,
		},
		{
			name: "Production environment with default config",
			config: config.Config{
				Environment: "production",
				Logging:     config.LoggingConfig{},
			},
			checkLevel:  slog.LevelInfo,
			checkType:   "json",
			checkOutput: os.Stdout,
		},
		{
			name: "Custom level configuration",
			config: config.Config{
				Environment: "production",
				Logging: config.LoggingConfig{
					Level: "ERROR",
				},
			},
			checkLevel:  slog.LevelError,
			checkType:   "json",
			checkOutput: os.Stdout,
		},
		{
			name: "Custom type configuration",
			config: config.Config{
				Environment: "development",
				Logging: config.LoggingConfig{
					Type: "JSON",
				},
			},
			checkLevel:  slog.LevelDebug,
			checkType:   "json",
			checkOutput: os.Stdout,
		},
		{
			name: "Output to stderr",
			config: config.Config{
				Environment: "production",
				Logging: config.LoggingConfig{
					Stderr: true,
				},
			},
			checkLevel:  slog.LevelInfo,
			checkType:   "json",
			checkOutput: os.Stderr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.config)

			// Check if logger is not null
			if logger == nil {
				t.Fatal("Logger should not be null")
			}

			// Get logger handler for verification
			handler := logger.Handler()

			// Check logger level
			if handler.Enabled(context.Background(), tt.checkLevel-1) {
				t.Errorf("Incorrect logger level. Expected: %v", tt.checkLevel)
			}

			// Check handler type
			switch tt.checkType {
			case "json":
				if _, ok := handler.(*slog.JSONHandler); !ok {
					t.Error("Handler should be of type JSONHandler")
				}
			case "text":
				if _, ok := handler.(*slog.TextHandler); !ok {
					t.Error("Handler should be of type TextHandler")
				}
			}
		})
	}
}
