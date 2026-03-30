package repository

import (
	"database/sql"

	"Phantom_backend/services/sticker-service/internal/models"
)

type PackRepository struct {
	db *sql.DB
}

func NewPackRepository(db *sql.DB) *PackRepository {
	return &PackRepository{db: db}
}

func (r *PackRepository) FindAll() ([]models.Pack, error) {
	rows, err := r.db.Query("SELECT id, name, description, preview_url, price, discount, created_at FROM sticker_packs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packs []models.Pack
	for rows.Next() {
		var p models.Pack
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PreviewURL,
			&p.Price, &p.Discount, &p.CreatedAt); err != nil {
			return nil, err
		}
		packs = append(packs, p)
	}

	return packs, nil
}
