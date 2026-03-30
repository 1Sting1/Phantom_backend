package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/forum-service/internal/models"
	"github.com/google/uuid"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) Create(thread *models.Thread) error {
	thread.ID = uuid.New().String()
	thread.CreatedAt = time.Now()
	thread.UpdatedAt = time.Now()

	query := `INSERT INTO threads (id, category_id, user_id, title, slug, content, is_pinned, is_locked, views_count, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(query, thread.ID, thread.CategoryID, thread.UserID, thread.Title,
		thread.Slug, thread.Content, thread.IsPinned, thread.IsLocked, thread.ViewsCount,
		thread.CreatedAt, thread.UpdatedAt)
	return err
}

func (r *ThreadRepository) FindAll() ([]models.Thread, error) {
	rows, err := r.db.Query("SELECT id, category_id, user_id, title, slug, content, is_pinned, is_locked, views_count, created_at, updated_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		if err := rows.Scan(&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Slug, &t.Content,
			&t.IsPinned, &t.IsLocked, &t.ViewsCount, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func (r *ThreadRepository) FindByID(id string) (*models.Thread, error) {
	thread := &models.Thread{}
	query := `SELECT id, category_id, user_id, title, slug, content, is_pinned, is_locked, views_count, created_at, updated_at 
	          FROM threads WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&thread.ID, &thread.CategoryID, &thread.UserID, &thread.Title,
		&thread.Slug, &thread.Content, &thread.IsPinned, &thread.IsLocked, &thread.ViewsCount,
		&thread.CreatedAt, &thread.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return thread, err
}

func (r *ThreadRepository) Update(thread *models.Thread) error {
	thread.UpdatedAt = time.Now()

	query := `UPDATE threads SET title = $1, content = $2, is_pinned = $3, is_locked = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(query, thread.Title, thread.Content, thread.IsPinned, thread.IsLocked, thread.UpdatedAt, thread.ID)
	return err
}

func (r *ThreadRepository) Delete(id string) error {
	query := `DELETE FROM threads WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ThreadRepository) IncrementViews(id string) error {
	query := `UPDATE threads SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
