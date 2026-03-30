package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"Phantom_backend/pkg/jwt"
	"Phantom_backend/pkg/security"
	"Phantom_backend/services/auth-service/internal/models"
	"Phantom_backend/services/auth-service/internal/repository"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

type AuthService struct {
	userRepo     *repository.UserRepository
	tokenRepo    *repository.TokenRepository
	tokenService *jwt.TokenService
	jwtSecret    string
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		tokenService: jwt.NewTokenService(jwtSecret),
		jwtSecret:    jwtSecret,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:         req.Email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		IsActive:      true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.issueTokenPair(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !security.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokenPair(user)
}

const (
	accessTokenDuration  = 24 * time.Hour
	refreshTokenDuration = 30 * 24 * time.Hour // 30 days for "remember me"
)

func (s *AuthService) issueTokenPair(user *models.User) (*models.AuthResponse, error) {
	token, err := s.tokenService.GenerateToken(user.ID, user.Email, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateToken(user.ID, user.Email, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(refreshTokenDuration)
	if err := s.tokenRepo.Create(user.ID, hashToken(refreshToken), expiresAt); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (*models.AuthResponse, error) {
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	oldHash := hashToken(refreshToken)
	userID, _, err := s.tokenRepo.FindValidByHash(oldHash)
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, ErrInvalidRefreshToken
	}

	resp, err := s.issueTokenPair(user)
	if err != nil {
		return nil, err
	}
	// Rotate: revoke old refresh only after new pair is issued (avoids losing the session if issue fails).
	_ = s.tokenRepo.RevokeByHash(oldHash)
	return resp, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.tokenRepo.RevokeByHash(hashToken(refreshToken))
}

func (s *AuthService) GetUser(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
