package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/landing-service/internal/models"
	"github.com/google/uuid"
)

type WaitlistRepository struct {
	db *sql.DB
}

func NewWaitlistRepository(db *sql.DB) *WaitlistRepository {
	return &WaitlistRepository{db: db}
}

func (r *WaitlistRepository) Create(entry *models.WaitlistEntry) error {
	entry.ID = uuid.New().String()
	entry.Status = "pending"
	entry.CreatedAt = time.Now()

	query := `INSERT INTO waitlist_entries (id, email, telegram, discord, status, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, entry.ID, entry.Email, entry.Telegram, entry.Discord, entry.Status, entry.CreatedAt)
	return err
}

func (r *WaitlistRepository) FindByEmail(email string) (*models.WaitlistEntry, error) {
	entry := &models.WaitlistEntry{}
	query := `SELECT id, email, telegram, discord, status, created_at 
	          FROM waitlist_entries WHERE email = $1`

	var telegram, discord sql.NullString
	err := r.db.QueryRow(query, email).Scan(
		&entry.ID, &entry.Email, &telegram, &discord, &entry.Status, &entry.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if telegram.Valid {
		entry.Telegram = &telegram.String
	}
	if discord.Valid {
		entry.Discord = &discord.String
	}

	return entry, nil
}

func (r *WaitlistRepository) FindByID(id string) (*models.WaitlistEntry, error) {
	entry := &models.WaitlistEntry{}
	query := `SELECT id, email, telegram, discord, status, created_at 
	          FROM waitlist_entries WHERE id = $1`

	var telegram, discord sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&entry.ID, &entry.Email, &telegram, &discord, &entry.Status, &entry.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if telegram.Valid {
		entry.Telegram = &telegram.String
	}
	if discord.Valid {
		entry.Discord = &discord.String
	}

	return entry, nil
}

func (r *WaitlistRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE waitlist_entries SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}
