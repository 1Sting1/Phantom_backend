package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/auth-service/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `INSERT INTO users (id, email, password_hash, created_at, updated_at, email_verified, is_active) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.EmailVerified, user.IsActive)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at, updated_at, email_verified, is_active 
	          FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt,
		&user.UpdatedAt, &user.EmailVerified, &user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at, updated_at, email_verified, is_active 
	          FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt,
		&user.UpdatedAt, &user.EmailVerified, &user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// TokenRepository handles refresh token storage
type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(userID, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked)
	          VALUES ($1, $2, $3, $4, $5, FALSE)`
	_, err := r.db.Exec(query, uuid.New().String(), userID, tokenHash, expiresAt, time.Now())
	return err
}

func (r *TokenRepository) FindValidByHash(tokenHash string) (userID string, expiresAt time.Time, err error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens 
	          WHERE token_hash = $1 AND revoked = FALSE AND expires_at > NOW()`
	err = r.db.QueryRow(query, tokenHash).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return "", time.Time{}, nil
	}
	return userID, expiresAt, err
}

func (r *TokenRepository) RevokeByHash(tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = $1`
	_, err := r.db.Exec(query, tokenHash)
	return err
}

func (r *TokenRepository) RevokeAllForUser(userID string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *TokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}
