// Package repository implementa domain.Repository usando bun/pgx como
// mecanismo de persistência. É a única camada que sabe falar SQL.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"go-api-base/internal/modules/user/domain"
	"go-api-base/internal/modules/user/model"
)

type userRepository struct {
	db *bun.DB
}

// NewUserRepository cria uma implementação de domain.Repository baseada em bun.
func NewUserRepository(db *bun.DB) domain.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	m := model.FromDomain(u)
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	if _, err := r.db.NewInsert().Model(m).Exec(ctx); err != nil {
		return err
	}

	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	m := new(model.User)

	err := r.db.NewSelect().Model(m).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m := new(model.User)

	err := r.db.NewSelect().Model(m).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	res, err := r.db.NewUpdate().
		Model((*model.User)(nil)).
		Set("name = ?", u.Name).
		Set("email = ?", u.Email).
		Set("updated_at = current_timestamp").
		Where("id = ?", u.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Como o modelo usa a tag "soft_delete", isso executa um UPDATE
	// definindo deleted_at, não um DELETE físico.
	res, err := r.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	var models []model.User

	count, err := r.db.NewSelect().
		Model(&models).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	users := make([]*domain.User, len(models))
	for i := range models {
		users[i] = models[i].ToDomain()
	}

	return users, count, nil
}
