package repository

import (
	"database/sql"

	"Phantom_backend/services/landing-service/internal/models"
)

type CarouselRepository struct {
	db *sql.DB
}

func NewCarouselRepository(db *sql.DB) *CarouselRepository {
	return &CarouselRepository{db: db}
}

func (r *CarouselRepository) FindAll() ([]models.CarouselSlide, error) {
	rows, err := r.db.Query("SELECT id, image_url, title, description, \"order\" FROM carousel_slides ORDER BY \"order\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slides []models.CarouselSlide
	for rows.Next() {
		var slide models.CarouselSlide
		var title, description sql.NullString
		if err := rows.Scan(&slide.ID, &slide.ImageURL, &title, &description, &slide.Order); err != nil {
			return nil, err
		}
		if title.Valid {
			slide.Title = &title.String
		}
		if description.Valid {
			slide.Description = &description.String
		}
		slides = append(slides, slide)
	}

	return slides, nil
}
