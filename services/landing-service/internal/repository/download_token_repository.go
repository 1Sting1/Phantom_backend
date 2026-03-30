package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/landing-service/internal/models"
	"github.com/google/uuid"
)

type DownloadTokenRepository struct {
	db *sql.DB
}

func NewDownloadTokenRepository(db *sql.DB) *DownloadTokenRepository {
	return &DownloadTokenRepository{db: db}
}

// GenerateToken creates a new download token
func (r *DownloadTokenRepository) GenerateToken(email *string, expiresInHours int) (*models.DownloadToken, error) {
	token := &models.DownloadToken{
		ID:        uuid.New().String(),
		Token:     uuid.New().String() + uuid.New().String(), // Generate a longer token
		Email:     email,
		ExpiresAt: time.Now().Add(time.Duration(expiresInHours) * time.Hour),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO download_tokens (id, token, email, expires_at, is_used, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, token.ID, token.Token, token.Email, token.ExpiresAt, token.IsUsed, token.CreatedAt)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// FindByToken finds a token by its value
func (r *DownloadTokenRepository) FindByToken(token string) (*models.DownloadToken, error) {
	dt := &models.DownloadToken{}
	query := `SELECT id, token, email, expires_at, is_used, created_at 
	          FROM download_tokens WHERE token = $1`

	var email sql.NullString
	err := r.db.QueryRow(query, token).Scan(
		&dt.ID, &dt.Token, &email, &dt.ExpiresAt, &dt.IsUsed, &dt.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if email.Valid {
		dt.Email = &email.String
	}

	return dt, nil
}

// MarkAsUsed marks a token as used
func (r *DownloadTokenRepository) MarkAsUsed(tokenID string) error {
	query := `UPDATE download_tokens SET is_used = TRUE WHERE id = $1`
	_, err := r.db.Exec(query, tokenID)
	return err
}

// IsValid checks if token is valid (not expired and not used)
func (r *DownloadTokenRepository) IsValid(token string) (bool, error) {
	dt, err := r.FindByToken(token)
	if err != nil {
		return false, err
	}

	if dt == nil {
		return false, nil
	}

	if dt.IsUsed {
		return false, nil
	}

	if time.Now().After(dt.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

// CleanupExpiredTokens removes expired tokens (can be called periodically)
func (r *DownloadTokenRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM download_tokens WHERE expires_at < NOW() OR is_used = TRUE`
	_, err := r.db.Exec(query)
	return err
}
