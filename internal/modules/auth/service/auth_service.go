// Package service implementa domain.Service do módulo auth. Depende das
// interfaces domain.Repository (auth) e domain.Repository (user), nunca
// de implementações concretas.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"go-api-base/internal/modules/auth/domain"
	userdomain "go-api-base/internal/modules/user/domain"
	"go-api-base/internal/observability"
	"go-api-base/pkg/hash"
	"go-api-base/pkg/jwt"
)

type authService struct {
	repo            domain.Repository
	userRepo        userdomain.Repository
	jwtManager      *jwt.Manager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewAuthService cria uma implementação de domain.Service.
func NewAuthService(
	repo domain.Repository,
	userRepo userdomain.Repository,
	jwtManager *jwt.Manager,
	accessTokenTTL, refreshTokenTTL time.Duration,
) domain.Service {
	return &authService{
		repo:            repo,
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*domain.TokenPair, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, userdomain.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, userdomain.ErrEmailAlreadyInUse
	}

	hashed, err := hash.Hash(password)
	if err != nil {
		return nil, err
	}

	u := &userdomain.User{Name: name, Email: email, Password: hashed}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, u)
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !hash.Compare(u.Password, password) {
		return nil, domain.ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	rt, err := s.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return nil, domain.ErrRefreshTokenExpired
	}

	u, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// Rotaciona o refresh token: o antigo é revogado e um novo é emitido.
	if err := s.repo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, u)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

// issueTokens gera um novo access token (JWT) e um novo refresh token
// (opaco, persistido no banco) para o usuário informado.
func (s *authService) issueTokens(ctx context.Context, u *userdomain.User) (*domain.TokenPair, error) {
	access, err := s.jwtManager.Generate(u.ID, u.Email, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refresh := uuid.NewString()
	rt := &domain.RefreshToken{
		UserID:    u.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}

	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	// Registra métricas
	observability.AuthTokensIssued.WithLabelValues("access").Inc()
	observability.AuthTokensIssued.WithLabelValues("refresh").Inc()

	return &domain.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}
