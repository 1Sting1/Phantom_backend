package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	pkgHttp "Phantom_backend/pkg/http"
	"Phantom_backend/services/landing-service/internal/models"
	"Phantom_backend/services/landing-service/internal/repository"
)

type LandingHandler struct {
	navRepo             *repository.NavigationRepository
	carouselRepo        *repository.CarouselRepository
	featureRepo         *repository.FeatureRepository
	waitlistRepo        *repository.WaitlistRepository
	waitlistConfirmRepo *repository.WaitlistConfirmationRepository
	releaseRepo         *repository.ReleaseRepository
	downloadTokenRepo   *repository.DownloadTokenRepository
	pageRepo            *repository.PageRepository
	footerRepo          *repository.FooterRepository
	communityRepo       *repository.CommunityRepository
}

func NewLandingHandler(db *sql.DB) *LandingHandler {
	return &LandingHandler{
		navRepo:             repository.NewNavigationRepository(db),
		carouselRepo:        repository.NewCarouselRepository(db),
		featureRepo:         repository.NewFeatureRepository(db),
		waitlistRepo:        repository.NewWaitlistRepository(db),
		waitlistConfirmRepo: repository.NewWaitlistConfirmationRepository(db),
		releaseRepo:         repository.NewReleaseRepository(db),
		downloadTokenRepo:   repository.NewDownloadTokenRepository(db),
		pageRepo:            repository.NewPageRepository(db),
		footerRepo:          repository.NewFooterRepository(db),
		communityRepo:       repository.NewCommunityRepository(db),
	}
}

// GetNavigation returns navigation items
func (h *LandingHandler) GetNavigation(w http.ResponseWriter, r *http.Request) {
	items, err := h.navRepo.FindAll()
	if err != nil {
		http.Error(w, "Failed to fetch navigation", http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, items)
}

// GetCarousel returns carousel slides
func (h *LandingHandler) GetCarousel(w http.ResponseWriter, r *http.Request) {
	slides, err := h.carouselRepo.FindAll()
	if err != nil {
		http.Error(w, "Failed to fetch carousel", http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, slides)
}

// GetFeatures returns features list
func (h *LandingHandler) GetFeatures(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	features, err := h.featureRepo.FindAll(language)
	if err != nil {
		http.Error(w, "Failed to fetch features", http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, features)
}

// AddToWaitlist adds an entry to the waitlist
func (h *LandingHandler) AddToWaitlist(w http.ResponseWriter, r *http.Request) {
	var req models.WaitlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Check if already exists
	existing, err := h.waitlistRepo.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, "Failed to check waitlist", http.StatusInternalServerError)
		return
	}

	if existing != nil {
		pkgHttp.Success(w, models.WaitlistResponse{
			ID:     existing.ID,
			Status: existing.Status,
		})
		return
	}

	entry := &models.WaitlistEntry{
		Email:    req.Email,
		Telegram: req.Telegram,
		Discord:  req.Discord,
	}

	if err := h.waitlistRepo.Create(entry); err != nil {
		http.Error(w, "Failed to add to waitlist", http.StatusInternalServerError)
		return
	}

	// Generate confirmation token (24 hours expiry)
	confirmationToken, err := h.waitlistConfirmRepo.GenerateToken(entry.ID, entry.Email, 24)
	if err != nil {
		// Log error but don't fail the request
		// In production, you might want to send email with confirmation link
	}

	response := models.WaitlistResponse{
		ID:     entry.ID,
		Status: entry.Status,
	}

	// Include confirmation token in response (in production, send via email)
	if confirmationToken != nil {
		// In real implementation, send email with confirmation link
		// For now, we'll just return success
	}

	pkgHttp.Created(w, response)
}

// ConfirmWaitlist confirms a waitlist entry using confirmation token
func (h *LandingHandler) ConfirmWaitlist(w http.ResponseWriter, r *http.Request) {
	var req models.WaitlistConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	err := h.waitlistConfirmRepo.ConfirmWaitlist(req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to confirm waitlist", http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, map[string]string{
		"message": "Waitlist entry confirmed successfully",
	})
}

// GetLatestRelease returns the latest release for a specific OS
func (h *LandingHandler) GetLatestRelease(w http.ResponseWriter, r *http.Request) {
	os := r.URL.Query().Get("os")
	if os == "" {
		http.Error(w, "OS parameter is required", http.StatusBadRequest)
		return
	}

	os = strings.ToLower(os)
	if os != "windows" && os != "linux" && os != "mac" {
		http.Error(w, "Invalid OS. Must be windows, linux, or mac", http.StatusBadRequest)
		return
	}

	release, err := h.releaseRepo.FindLatest(os)
	if err != nil {
		http.Error(w, "Failed to fetch release", http.StatusInternalServerError)
		return
	}

	if release == nil {
		http.Error(w, "No release found", http.StatusNotFound)
		return
	}

	pkgHttp.Success(w, release)
}

// GetPage returns a page by slug
func (h *LandingHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "Slug parameter is required", http.StatusBadRequest)
		return
	}

	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "en"
	}

	page, err := h.pageRepo.FindBySlug(slug, language)
	if err != nil {
		http.Error(w, "Failed to fetch page", http.StatusInternalServerError)
		return
	}

	if page == nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	pkgHttp.Success(w, page)
}

// GetFooterLinks returns footer links
func (h *LandingHandler) GetFooterLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.footerRepo.FindAll()
	if err != nil {
		http.Error(w, "Failed to fetch footer links", http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, links)
}

// GetCommunityLink returns active community link
func (h *LandingHandler) GetCommunityLink(w http.ResponseWriter, r *http.Request) {
	link, err := h.communityRepo.FindActive()
	if err != nil {
		http.Error(w, "Failed to fetch community link", http.StatusInternalServerError)
		return
	}

	if link == nil {
		http.Error(w, "No active community link found", http.StatusNotFound)
		return
	}

	pkgHttp.Success(w, link)
}

// GetDownloadToken generates a download token (optional: requires waitlist approval)
func (h *LandingHandler) GetDownloadToken(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}

	// Generate token with 1 hour expiry
	token, err := h.downloadTokenRepo.GenerateToken(emailPtr, 1)
	if err != nil {
		http.Error(w, "Failed to generate download token", http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, models.DownloadTokenResponse{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	})
}

// Download handles file download with optional token validation
func (h *LandingHandler) Download(w http.ResponseWriter, r *http.Request) {
	os := r.URL.Query().Get("os")
	token := r.URL.Query().Get("token")

	if os == "" {
		http.Error(w, "OS parameter is required", http.StatusBadRequest)
		return
	}

	os = strings.ToLower(os)
	if os != "windows" && os != "linux" && os != "mac" {
		http.Error(w, "Invalid OS. Must be windows, linux, or mac", http.StatusBadRequest)
		return
	}

	// If token is provided, validate it
	if token != "" {
		isValid, err := h.downloadTokenRepo.IsValid(token)
		if err != nil {
			http.Error(w, "Failed to validate token", http.StatusInternalServerError)
			return
		}
		if !isValid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Mark token as used
		dt, err := h.downloadTokenRepo.FindByToken(token)
		if err == nil && dt != nil {
			h.downloadTokenRepo.MarkAsUsed(dt.ID)
		}
	}

	// Get latest release
	release, err := h.releaseRepo.FindLatest(os)
	if err != nil {
		http.Error(w, "Failed to fetch release", http.StatusInternalServerError)
		return
	}

	if release == nil {
		http.Error(w, "No release found for this OS", http.StatusNotFound)
		return
	}

	// Redirect to download URL
	http.Redirect(w, r, release.DownloadURL, http.StatusFound)
}

// HealthCheck returns service health status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.Success(w, map[string]string{
		"status":  "ok",
		"service": "landing-service",
	})
}
