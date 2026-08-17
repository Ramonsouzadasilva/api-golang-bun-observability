package seeds

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"go-api-base/internal/modules/user/domain"
	"go-api-base/internal/modules/user/repository"
	"go-api-base/pkg/hash"
)

// Run popula o banco com dados iniciais necessários para começar a usar
// a API (hoje, apenas um usuário administrador de exemplo).
func Run(ctx context.Context, db *bun.DB) error {
	repo := repository.NewUserRepository(db)

	_, err := repo.FindByEmail(ctx, "admin@example.com")
	if err == nil {
		fmt.Println("admin user already exists, skipping seed")
		return nil
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	hashed, err := hash.Hash("Admin@123")
	if err != nil {
		return err
	}

	admin := &domain.User{
		Name:     "Administrator",
		Email:    "admin@example.com",
		Password: hashed,
	}

	if err := repo.Create(ctx, admin); err != nil {
		return err
	}

	fmt.Println("admin user created: admin@example.com / Admin@123")
	return nil
}
