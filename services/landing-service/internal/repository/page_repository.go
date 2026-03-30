package repository

import (
	"database/sql"

	"Phantom_backend/services/landing-service/internal/models"
)

type PageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) FindBySlug(slug, language string) (*models.Page, error) {
	page := &models.Page{}
	query := `SELECT id, slug, title, content, language, created_at, updated_at 
	          FROM pages WHERE slug = $1 AND language = $2`

	err := r.db.QueryRow(query, slug, language).Scan(
		&page.ID, &page.Slug, &page.Title, &page.Content,
		&page.Language, &page.CreatedAt, &page.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return page, err
}
