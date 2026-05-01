package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Env         string
	SeedOnEmpty bool
}

func Load() *Config {
	for _, p := range []string{".env", "../.env", "../../.env", "../../../.env"} {
		if err := godotenv.Load(p); err == nil {
			break
		}
	}

	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://pharmasense:pharmasense@127.0.0.1:5433/pharmasense?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-key-change-in-production-32c"),
		Port:        getEnv("PORT", "3001"),
		Env:         getEnv("ENV", "development"),
		SeedOnEmpty: getEnv("SEED_ON_EMPTY", "true") == "true",
	}

	slog.Info("config loaded", "port", cfg.Port, "env", cfg.Env)
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
