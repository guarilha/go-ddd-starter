package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/guarilha/go-ddd-starter/app/cli/user"
	"github.com/guarilha/go-ddd-starter/domain"
	"github.com/guarilha/go-ddd-starter/internal/config"
	"github.com/guarilha/go-ddd-starter/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v3"
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

	// Domains
	// ------------------------------------------
	exampleConfig := domain.Config{
		Example: "example",
	}

	domains, err := domain.NewDomains(conn, exampleConfig)
	if err != nil {
		mainLogger.Error("failed to create domains", "error", err)
		return
	}

	// CLI
	// ------------------------------------------

	newUser := user.NewUserCommand{}

	cmd := &cli.Command{
		Name:  "cli",
		Usage: "Manage your local Go DDD Starter Instance",
		Commands: []*cli.Command{
			{
				Name:  "user",
				Usage: "Create users",
				Commands: []*cli.Command{
					newUser.Command(domains),
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
