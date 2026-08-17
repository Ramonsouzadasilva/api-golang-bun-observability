// Package domain define as entidades e contratos do módulo de autenticação.
package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired or revoked")
)

// RefreshToken representa um token de renovação persistido, usado para
// obter novos access tokens sem exigir login novamente.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

// TokenPair é o par de tokens devolvido ao cliente após login/registro/refresh.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Repository define a persistência de refresh tokens.
type Repository interface {
	CreateRefreshToken(ctx context.Context, rt *RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}

// Service define as regras de negócio de autenticação.
type Service interface {
	Register(ctx context.Context, name, email, password string) (*TokenPair, error)
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}
