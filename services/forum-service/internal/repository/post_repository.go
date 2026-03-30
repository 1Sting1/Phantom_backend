package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/forum-service/internal/models"
	"github.com/google/uuid"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()

	var parentPostID sql.NullString
	if post.ParentPostID != "" {
		parentPostID = sql.NullString{String: post.ParentPostID, Valid: true}
	}

	query := `INSERT INTO posts (id, thread_id, user_id, parent_post_id, content, is_edited, edited_at, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query, post.ID, post.ThreadID, post.UserID, parentPostID,
		post.Content, post.IsEdited, post.EditedAt, post.CreatedAt)
	return err
}

func (r *PostRepository) FindByThreadID(threadID string) ([]models.Post, error) {
	rows, err := r.db.Query("SELECT id, thread_id, user_id, parent_post_id, content, is_edited, edited_at, created_at FROM posts WHERE thread_id = $1 ORDER BY created_at ASC", threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var parentPostID sql.NullString
		var editedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.UserID, &parentPostID, &p.Content, &p.IsEdited, &editedAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		if parentPostID.Valid {
			p.ParentPostID = parentPostID.String
		}
		if editedAt.Valid {
			p.EditedAt = editedAt.Time
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id string) (*models.Post, error) {
	post := &models.Post{}
	query := `SELECT id, thread_id, user_id, parent_post_id, content, is_edited, edited_at, created_at 
	          FROM posts WHERE id = $1`

	var parentPostID sql.NullString
	var editedAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.ThreadID, &post.UserID, &parentPostID,
		&post.Content, &post.IsEdited, &editedAt, &post.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if parentPostID.Valid {
		post.ParentPostID = parentPostID.String
	}
	if editedAt.Valid {
		post.EditedAt = editedAt.Time
	}

	return post, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	post.IsEdited = true
	post.EditedAt = time.Now()

	query := `UPDATE posts SET content = $1, is_edited = $2, edited_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, post.Content, post.IsEdited, post.EditedAt, post.ID)
	return err
}

func (r *PostRepository) Delete(id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
