package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	pkgHttp "Phantom_backend/pkg/http"
	"Phantom_backend/services/forum-service/internal/models"
	"Phantom_backend/services/forum-service/internal/repository"

	"github.com/gorilla/mux"
)

type ForumHandler struct {
	categoryRepo *repository.CategoryRepository
	threadRepo   *repository.ThreadRepository
	postRepo     *repository.PostRepository
}

func NewForumHandler(db *sql.DB) *ForumHandler {
	return &ForumHandler{
		categoryRepo: repository.NewCategoryRepository(db),
		threadRepo:   repository.NewThreadRepository(db),
		postRepo:     repository.NewPostRepository(db),
	}
}

func (h *ForumHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkgHttp.Success(w, categories)
}

func (h *ForumHandler) GetThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := h.threadRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkgHttp.Success(w, threads)
}

func (h *ForumHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID := vars["id"]

	thread, err := h.threadRepo.FindByID(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if thread == nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	pkgHttp.Success(w, thread)
}

func (h *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Authorization required to create a thread", http.StatusUnauthorized)
		return
	}

	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if thread.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if thread.CategoryID == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	thread.UserID = userID
	if thread.Slug == "" {
		thread.Slug = thread.Title // simple slug; could be sanitized
	}
	if err := h.threadRepo.Create(&thread); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Created(w, thread)
}

func (h *ForumHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		// Anonymous replies (no Bearer): attribute to well-known guest id so clients without auth still work.
		userID = models.GuestUserID
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.ThreadID == "" || post.Content == "" {
		http.Error(w, "Thread ID and content are required", http.StatusBadRequest)
		return
	}

	post.UserID = userID
	if err := h.postRepo.Create(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Created(w, post)
}

func (h *ForumHandler) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID := vars["id"]

	posts, err := h.postRepo.FindByThreadID(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, posts)
}

func (h *ForumHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Get existing thread
	thread, err := h.threadRepo.FindByID(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if thread == nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	if userID == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	if thread.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var updatedThread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&updatedThread); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	thread.Title = updatedThread.Title
	thread.Content = updatedThread.Content
	if updatedThread.IsPinned {
		thread.IsPinned = updatedThread.IsPinned
	}
	if updatedThread.IsLocked {
		thread.IsLocked = updatedThread.IsLocked
	}

	if err := h.threadRepo.Update(thread); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, thread)
}

func (h *ForumHandler) DeleteThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Get existing thread
	thread, err := h.threadRepo.FindByID(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if thread == nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	if userID == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	if thread.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.threadRepo.Delete(threadID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, map[string]string{"message": "Thread deleted successfully"})
}

func (h *ForumHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Get existing post
	post, err := h.postRepo.FindByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if userID == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	if post.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post.Content = updatedPost.Content
	if err := h.postRepo.Update(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, post)
}

func (h *ForumHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Get existing post
	post, err := h.postRepo.FindByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if userID == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	if post.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.postRepo.Delete(postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, map[string]string{"message": "Post deleted successfully"})
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.Success(w, map[string]string{
		"status":  "ok",
		"service": "forum-service",
	})
}
