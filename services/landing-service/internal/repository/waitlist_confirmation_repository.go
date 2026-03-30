package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/landing-service/internal/models"
	"github.com/google/uuid"
)

type WaitlistConfirmationRepository struct {
	db *sql.DB
}

func NewWaitlistConfirmationRepository(db *sql.DB) *WaitlistConfirmationRepository {
	return &WaitlistConfirmationRepository{db: db}
}

// GenerateToken creates a new confirmation token for a waitlist entry
func (r *WaitlistConfirmationRepository) GenerateToken(waitlistID, email string, expiresInHours int) (*models.WaitlistConfirmationToken, error) {
	token := &models.WaitlistConfirmationToken{
		ID:         uuid.New().String(),
		Token:      uuid.New().String() + uuid.New().String(),
		WaitlistID: waitlistID,
		Email:      email,
		ExpiresAt:  time.Now().Add(time.Duration(expiresInHours) * time.Hour),
		IsUsed:     false,
		CreatedAt:  time.Now(),
	}

	query := `INSERT INTO waitlist_confirmation_tokens (id, token, waitlist_id, email, expires_at, is_used, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query, token.ID, token.Token, token.WaitlistID, token.Email,
		token.ExpiresAt, token.IsUsed, token.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Update waitlist entry with confirmation token ID
	updateQuery := `UPDATE waitlist_entries SET confirmation_token_id = $1 WHERE id = $2`
	_, err = r.db.Exec(updateQuery, token.ID, waitlistID)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// FindByToken finds a confirmation token by its value
func (r *WaitlistConfirmationRepository) FindByToken(token string) (*models.WaitlistConfirmationToken, error) {
	wct := &models.WaitlistConfirmationToken{}
	query := `SELECT id, token, waitlist_id, email, expires_at, is_used, created_at 
	          FROM waitlist_confirmation_tokens WHERE token = $1`

	err := r.db.QueryRow(query, token).Scan(
		&wct.ID, &wct.Token, &wct.WaitlistID, &wct.Email,
		&wct.ExpiresAt, &wct.IsUsed, &wct.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return wct, nil
}

// ConfirmWaitlist marks token as used and updates waitlist entry status
func (r *WaitlistConfirmationRepository) ConfirmWaitlist(token string) error {
	// Find token
	wct, err := r.FindByToken(token)
	if err != nil {
		return err
	}

	if wct == nil {
		return sql.ErrNoRows
	}

	if wct.IsUsed {
		return sql.ErrNoRows // Already used
	}

	if time.Now().After(wct.ExpiresAt) {
		return sql.ErrNoRows // Expired
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Mark token as used
	updateTokenQuery := `UPDATE waitlist_confirmation_tokens SET is_used = TRUE WHERE id = $1`
	_, err = tx.Exec(updateTokenQuery, wct.ID)
	if err != nil {
		return err
	}

	// Update waitlist entry status and confirmed_at
	updateWaitlistQuery := `UPDATE waitlist_entries SET status = 'approved', confirmed_at = NOW() WHERE id = $1`
	_, err = tx.Exec(updateWaitlistQuery, wct.WaitlistID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// IsValid checks if token is valid
func (r *WaitlistConfirmationRepository) IsValid(token string) (bool, error) {
	wct, err := r.FindByToken(token)
	if err != nil {
		return false, err
	}

	if wct == nil {
		return false, nil
	}

	if wct.IsUsed {
		return false, nil
	}

	if time.Now().After(wct.ExpiresAt) {
		return false, nil
	}

	return true, nil
}
