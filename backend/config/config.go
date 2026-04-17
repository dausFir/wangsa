package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort   string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
	CookieDomain string
	IsProduction bool
	FrontendURL  string
}

func Load() *Config {
	// ACCESS_TOKEN_MINUTES: short-lived access token TTL (default 15 min)
	// REFRESH_TOKEN_DAYS: handled in jwt.Manager (hardcoded 30 days)
	accessMinutes, _ := strconv.Atoi(getEnv("ACCESS_TOKEN_MINUTES", "15"))
	isProd, _       := strconv.ParseBool(getEnv("PRODUCTION", "false"))

	return &Config{
		ServerPort:   getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://wangsa:wangsa@localhost:5432/wangsa?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "super-secret-jwt-key-change-in-production-please"),
		JWTExpiresIn: time.Duration(accessMinutes) * time.Minute,
		CookieDomain: getEnv("COOKIE_DOMAIN", "localhost"),
		IsProduction: isProd,
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

// ValidateDatabaseURL checks the DATABASE_URL is a parseable postgres:// URL.
func (c *Config) ValidateDatabaseURL() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		return fmt.Errorf("DATABASE_URL is not a valid URL: %w", err)
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("DATABASE_URL must start with postgres:// or postgresql://, got: %s://", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("DATABASE_URL missing host")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
