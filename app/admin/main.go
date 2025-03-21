package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/guarilha/go-ddd-starter/app/service/api"
	"github.com/guarilha/go-ddd-starter/internal/config"
	"github.com/guarilha/go-ddd-starter/internal/logger"
)

// Injected on build time by ldflags.
var (
	BuildCommit = "undefined"
	BuildTime   = "undefined"
)

func main() {
	ctx := context.Background()

	var cfg config.Config
	if err := cfg.Load(); err != nil {
		panic(fmt.Errorf("loading config: %w", err))
	}

	// Logger (Slog)
	// ------------------------------------------
	defaultLogger := logger.NewLogger(cfg)
	mainLogger := defaultLogger.With(
		"environment", cfg.Environment,
		"build_commit", BuildCommit,
		"build_time", BuildTime,
		"go_max_procs", runtime.GOMAXPROCS(0),
		"runtime_num_cpu", runtime.NumCPU(),
	)

	router := api.Router(nil) // We're not using domains for now

	// SERVER
	// ------------------------------------------
	server := http.Server{
		Handler:           router,
		Addr:              cfg.AdminApiAddress,
		ReadHeaderTimeout: 60 * time.Second,
	}
	mainLogger.Info("server started",
		"address", server.Addr,
	)

	if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		mainLogger.Error("failed to listen and serve server",
			"error", serverErr,
		)
	}
}
