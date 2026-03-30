package config

import (
	"database/sql"
	"os"
	"strings"

	"Phantom_backend/pkg/database"
)

type Config struct {
	Port string
	DB   *sql.DB
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8003"),
	}
}

func (c *Config) InitDB() (*sql.DB, error) {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "forum_db"),
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

// RunMigrations applies forum schema and seed (idempotent). Used on startup so Docker does not require a separate migrate step.
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
	`CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) UNIQUE NOT NULL,
    parent_id VARCHAR(255),
    "order" INTEGER DEFAULT 0,
    icon VARCHAR(100),
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS threads (
    id VARCHAR(255) PRIMARY KEY,
    category_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500),
    content TEXT NOT NULL,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    views_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS posts (
    id VARCHAR(255) PRIMARY KEY,
    thread_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    parent_post_id VARCHAR(255),
    content TEXT NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_threads_category_id ON threads(category_id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts(thread_id);`,
	`INSERT INTO categories (id, name, description, slug, parent_id, "order", icon, created_at)
VALUES
  ('general', 'General Discussion', 'General topics and community chat', 'general', NULL, 0, '', NOW()),
  ('installation', 'Installation & Setup', 'Installation issues and setup help', 'installation', NULL, 1, '', NOW()),
  ('features', 'Feature Requests', 'Suggestions and feature requests', 'features', NULL, 2, '', NOW()),
  ('development', 'Development', 'Development and contributing', 'development', NULL, 3, '', NOW())
ON CONFLICT (id) DO NOTHING;`,
}
