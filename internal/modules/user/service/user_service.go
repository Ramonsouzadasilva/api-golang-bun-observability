// Package service implementa domain.Service, contendo as regras de
// negócio do módulo user. Depende apenas da interface domain.Repository,
// nunca de uma implementação concreta (Dependency Inversion).
package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"go-api-base/internal/modules/user/domain"
	"go-api-base/pkg/hash"
	"go-api-base/pkg/pagination"
)

type userService struct {
	repo domain.Repository
}

// NewUserService cria uma implementação de domain.Service.
func NewUserService(repo domain.Repository) domain.Service {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, name, email, password string) (*domain.User, error) {
	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	hashed, err := hash.Hash(password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{Name: name, Email: email, Password: hashed}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, name string) (*domain.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Name = name

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, page, perPage int) ([]*domain.User, int, error) {
	limit, offset := pagination.LimitOffset(page, perPage)
	return s.repo.List(ctx, limit, offset)
}
