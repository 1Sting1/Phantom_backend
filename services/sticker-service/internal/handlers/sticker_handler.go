package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httputil "Phantom_backend/pkg/http"
	"Phantom_backend/services/sticker-service/internal/models"
	"Phantom_backend/services/sticker-service/internal/repository"
)

type StickerHandler struct {
	stickerRepo *repository.StickerRepository
	packRepo    *repository.PackRepository
	orderRepo   *repository.OrderRepository
}

func NewStickerHandler(db *sql.DB) *StickerHandler {
	return &StickerHandler{
		stickerRepo: repository.NewStickerRepository(db),
		packRepo:    repository.NewPackRepository(db),
		orderRepo:   repository.NewOrderRepository(db),
	}
}

func (h *StickerHandler) GetStickers(w http.ResponseWriter, r *http.Request) {
	stickers, err := h.stickerRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputil.Success(w, stickers)
}

func (h *StickerHandler) GetPacks(w http.ResponseWriter, r *http.Request) {
	packs, err := h.packRepo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputil.Success(w, packs)
}

func (h *StickerHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order.ID = newOrderID()
	order.UserID = userID
	order.Status = "pending"
	if order.TotalAmount < 0 {
		order.TotalAmount = 0
	}

	if err := h.orderRepo.Create(&order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.Created(w, order)
}

func newOrderID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("ord-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, map[string]string{
		"status":  "ok",
		"service": "sticker-service",
	})
}
