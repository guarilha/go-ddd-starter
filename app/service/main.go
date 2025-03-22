package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/guarilha/go-ddd-starter/app/service/api"
	v1 "github.com/guarilha/go-ddd-starter/app/service/api/v1"
	"github.com/guarilha/go-ddd-starter/domain"
	"github.com/guarilha/go-ddd-starter/internal/config"
	"github.com/guarilha/go-ddd-starter/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
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

	// Repositories
	// ------------------------------------------
	conn, err := pgxpool.New(ctx, cfg.ConnectionString())
	if err != nil {
		mainLogger.Error("failed to setup postgres", "error", err)
		return
	}
	defer conn.Close()

	// Handlers V1 and their dependencies
	// ------------------------------------------
	exampleConfig := domain.Config{
		Example: "example",
	}

	domains, err := domain.NewDomains(conn, exampleConfig)
	if err != nil {
		mainLogger.Error("failed to create domains", "error", err)
		return
	}

	apiV1 := v1.ApiHandlers{
		UserDomain: domains.User,
	}

	router := api.Router()
	apiV1.Routes(router)

	// SERVER
	// ------------------------------------------
	server := http.Server{
		Handler:           router,
		Addr:              cfg.ApiAddress,
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
