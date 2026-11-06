// Package config loads and validates service configuration from the
// environment. Configuration is read exactly once, at startup.
package config

import "fmt"

// Config holds all runtime configuration for the API.
type Config struct {
	Addr        string // HTTP listen address, e.g. ":8080"
	DatabaseURL string // PostgreSQL DSN
	FrontendURL string // Allowed CORS origin (the Next.js dev server)
}

// Load builds a Config from the environment and validates it in one place,
// reporting every missing value at startup instead of failing later at use.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		Addr:        getenv("ADDR"),
		DatabaseURL: getenv("DATABASE_URL"),
		FrontendURL: getenv("FRONTEND_URL"),
	}

	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.FrontendURL == "" {
		cfg.FrontendURL = "http://localhost:3000"
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing environment variables: %v", missing)
	}
	return cfg, nil
}
