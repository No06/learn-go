package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"learn-go/internal/config"
	"learn-go/internal/domain"
	"learn-go/internal/repository"
	"learn-go/pkg/crypto"
)

// AuthService handles authentication and token issuance.
type AuthService struct {
	accounts repository.AccountRepository
	cfg      config.AppConfig
}

// NewAuthService creates a new AuthService.
func NewAuthService(accounts repository.AccountRepository, cfg config.AppConfig) *AuthService {
	return &AuthService{accounts: accounts, cfg: cfg}
}

// Login authenticates a user and returns JWT access and refresh tokens.
func (s *AuthService) Login(ctx context.Context, schoolID, identifier, password string) (string, string, *domain.Account, error) {
	account, err := s.accounts.FindByIdentifier(ctx, schoolID, identifier)
	if err != nil {
		return "", "", nil, err
	}
	if err := crypto.ComparePassword(account.PasswordHash, password); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateToken(account.ID, string(account.Role), s.cfg.JWTSecret, s.cfg.TokenTTL)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := s.generateToken(account.ID, string(account.Role), s.cfg.RefreshSecret, s.cfg.RefreshTokenTTL)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, account, nil
}

func (s *AuthService) generateToken(subjectID, role, secret string, ttlSeconds int64) (string, error) {
	claims := jwt.MapClaims{
		"sub":  subjectID,
		"role": role,
		"exp":  time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix(),
		"iat":  time.Now().Unix(),
		"jti":  uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
