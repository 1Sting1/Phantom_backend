package repository

import (
	"database/sql"

	"Phantom_backend/services/landing-service/internal/models"
)

type FooterRepository struct {
	db *sql.DB
}

func NewFooterRepository(db *sql.DB) *FooterRepository {
	return &FooterRepository{db: db}
}

func (r *FooterRepository) FindAll() ([]models.FooterLink, error) {
	rows, err := r.db.Query("SELECT id, title, href, category, \"order\" FROM footer_links ORDER BY category, \"order\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.FooterLink
	for rows.Next() {
		var link models.FooterLink
		if err := rows.Scan(&link.ID, &link.Title, &link.Href, &link.Category, &link.Order); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}
