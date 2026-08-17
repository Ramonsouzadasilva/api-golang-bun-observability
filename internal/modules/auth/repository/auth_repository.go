package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"go-api-base/internal/modules/auth/domain"
	"go-api-base/internal/modules/auth/model"
)

type authRepository struct {
	db *bun.DB
}

// NewAuthRepository cria uma implementação de domain.Repository baseada em bun.
func NewAuthRepository(db *bun.DB) domain.Repository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	m := model.FromDomain(rt)
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	if _, err := r.db.NewInsert().Model(m).Exec(ctx); err != nil {
		return err
	}

	rt.ID = m.ID
	rt.CreatedAt = m.CreatedAt
	return nil
}

func (r *authRepository) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	m := new(model.RefreshToken)

	err := r.db.NewSelect().Model(m).Where("token = ?", token).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	now := time.Now()

	_, err := r.db.NewUpdate().
		Model((*model.RefreshToken)(nil)).
		Set("revoked_at = ?", now).
		Where("token = ?", token).
		Exec(ctx)

	return err
}

func (r *authRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()

	_, err := r.db.NewUpdate().
		Model((*model.RefreshToken)(nil)).
		Set("revoked_at = ?", now).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Exec(ctx)

	return err
}
