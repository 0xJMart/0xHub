package app

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultHTTPAddr           = ":8080"
	defaultDatabaseURL        = "postgres://hub:hub@localhost:5432/hub?sslmode=disable"
	defaultShutdownTimeout    = 10 * time.Second
	defaultDatabaseMaxConns   = int32(10)
	defaultDatabaseMinConns   = int32(1)
	defaultDatabaseConnMaxAge = 30 * time.Minute
)

// Config captures runtime configuration for the application.
type Config struct {
	HTTPAddr                string
	DatabaseURL             string
	ShutdownTimeout         time.Duration
	DatabaseMaxConns        int32
	DatabaseMinConns        int32
	DatabaseMaxConnLifetime time.Duration
}

// LoadConfig derives configuration from environment variables with reasonable defaults.
func LoadConfig() Config {
	cfg := Config{
		HTTPAddr:                defaultHTTPAddr,
		DatabaseURL:             defaultDatabaseURL,
		ShutdownTimeout:         defaultShutdownTimeout,
		DatabaseMaxConns:        defaultDatabaseMaxConns,
		DatabaseMinConns:        defaultDatabaseMinConns,
		DatabaseMaxConnLifetime: defaultDatabaseConnMaxAge,
	}

	if v := os.Getenv("HUB_API_HTTP_ADDR"); v != "" {
		cfg.HTTPAddr = v
	}
	if v := os.Getenv("HUB_API_DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if v := os.Getenv("HUB_API_SHUTDOWN_TIMEOUT"); v != "" {
		if dur, err := time.ParseDuration(v); err == nil && dur > 0 {
			cfg.ShutdownTimeout = dur
		}
	}
	if v := os.Getenv("HUB_API_DB_MAX_CONNS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			cfg.DatabaseMaxConns = int32(parsed)
		}
	}
	if v := os.Getenv("HUB_API_DB_MIN_CONNS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			cfg.DatabaseMinConns = int32(parsed)
		}
	}
	if v := os.Getenv("HUB_API_DB_CONN_MAX_LIFETIME"); v != "" {
		if dur, err := time.ParseDuration(v); err == nil && dur > 0 {
			cfg.DatabaseMaxConnLifetime = dur
		}
	}

	return cfg
}
