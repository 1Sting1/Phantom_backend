package config

import (
	"database/sql"
	"os"
	"strings"

	"Phantom_backend/pkg/database"
)

type Config struct {
	Port      string
	JWTSecret string
	DB        *sql.DB
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8001"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func (c *Config) InitDB() (*sql.DB, error) {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "auth_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	c.DB = db
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// RunMigrations runs auth DB migrations (idempotent). Kept here so Docker build does not depend on internal/migrate.
func RunMigrations(db *sql.DB) error {
	for _, block := range migrationBlocks {
		for _, s := range splitSQL(block) {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			if _, err := db.Exec(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func splitSQL(block string) []string {
	var out []string
	for _, s := range strings.Split(block, ";") {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t+";")
		}
	}
	return out
}

var migrationBlocks = []string{
	`CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
	`CREATE TABLE IF NOT EXISTS refresh_tokens (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);`,
}
