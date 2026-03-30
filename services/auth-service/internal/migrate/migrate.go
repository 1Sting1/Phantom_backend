package migrate

import (
	"database/sql"
	"strings"
)

// Run executes auth-service migrations in order (idempotent: SQL uses IF NOT EXISTS).
func Run(db *sql.DB) error {
	for _, block := range migrations {
		for _, s := range splitStatements(block) {
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

func splitStatements(block string) []string {
	var out []string
	for _, s := range strings.Split(block, ";") {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t+";")
		}
	}
	return out
}

var migrations = []string{
	migration001,
	migration002,
}

const migration001 = `
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`

const migration002 = `
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
`
