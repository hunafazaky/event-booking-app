package config

import (
	"fmt"
	"os"
)

// Config holds every environment-derived setting the app needs, loaded once
// at startup. Passing this struct down (instead of calling os.Getenv all
// over the codebase) makes it obvious what each package actually depends on,
// and makes it possible to construct a Config by hand in tests.
type Config struct {
	Port string

	DBURI        string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresSSL  string

	JWTSecret          string
	ImageKitPrivateKey string
}

// Load reads environment variables into a Config and validates that
// everything required is present. Call this once in main() — if it fails,
// fail fast at startup rather than hundreds of requests in.
func Load() (*Config, error) {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),

		DBURI:        os.Getenv("DB_URI"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPass: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:   os.Getenv("POSTGRES_DB_NAME"),
		PostgresSSL:  os.Getenv("POSTGRES_SSLMODE"),

		JWTSecret:          os.Getenv("JWT_SECRET"),
		ImageKitPrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
	}

	if cfg.DBURI == "" {
		missing := []string{}
		if cfg.PostgresHost == "" {
			missing = append(missing, "POSTGRES_HOST")
		}
		if cfg.PostgresPort == "" {
			missing = append(missing, "POSTGRES_PORT")
		}
		if cfg.PostgresUser == "" {
			missing = append(missing, "POSTGRES_USER")
		}
		if cfg.PostgresPass == "" {
			missing = append(missing, "POSTGRES_PASSWORD")
		}
		if cfg.PostgresDB == "" {
			missing = append(missing, "POSTGRES_DB_NAME")
		}
		if cfg.PostgresSSL == "" {
			missing = append(missing, "POSTGRES_SSLMODE")
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("missing required environment variables (or set DB_URI instead): %v", missing)
		}
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// DSN builds the Postgres connection string, preferring DB_URI when set.
func (c *Config) DSN() string {
	if c.DBURI != "" {
		return c.DBURI
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.PostgresHost, c.PostgresUser, c.PostgresPass, c.PostgresDB, c.PostgresPort, c.PostgresSSL,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
