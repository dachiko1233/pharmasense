package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"pharmasense/backend/cmd/seed"
	"pharmasense/backend/internal/config"
	embeddedmigrations "pharmasense/backend/internal/migrations"
	"pharmasense/backend/internal/repository"
	"pharmasense/backend/internal/server"
	"pharmasense/backend/internal/services"
)

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "run migrations and exit")
	seedOnly := flag.Bool("seed-only", false, "run seed and exit")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.Load()
	ctx := context.Background()

	pool, err := repository.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("db connect failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	migSQL, err := loadMigration()
	if err != nil {
		slog.Error("load migration", "error", err)
		os.Exit(1)
	}
	if err := repository.RunMigrations(ctx, pool, migSQL); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	if *migrateOnly {
		slog.Info("migrations complete")
		return
	}

	empty, err := repository.IsEmpty(ctx, pool)
	if err != nil {
		slog.Warn("could not check db emptiness", "error", err)
	}
	if (empty && cfg.SeedOnEmpty) || *seedOnly {
		slog.Info("seeding demo data...")
		if err := seed.RunSeed(ctx, pool); err != nil {
			slog.Error("seed failed", "error", err)
			if *seedOnly {
				os.Exit(1)
			}
		}
	}

	if *seedOnly {
		return
	}

	repo := repository.New(pool)
	authSvc := services.NewAuthService(repo, cfg.JWTSecret)
	engine := services.NewRiskEngine(repo)

	h := server.New(repo, authSvc, engine)

	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("server starting", "addr", addr)

	if err := http.ListenAndServe(addr, h); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func loadMigration() (string, error) {
	entries, err := embeddedmigrations.FS.ReadDir(".")
	if err != nil {
		return "", fmt.Errorf("read migrations dir: %w", err)
	}

	var parts []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			data, err := embeddedmigrations.FS.ReadFile(e.Name())
			if err != nil {
				return "", err
			}
			parts = append(parts, string(data))
		}
	}
	return strings.Join(parts, "\n"), nil
}
