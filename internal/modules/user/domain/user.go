// Package domain define a entidade de negócio User e os contratos
// (interfaces) que as camadas de service e repository implementam.
// Nenhuma outra camada deve depender de detalhes de HTTP ou banco aqui,
// isso é o que garante a inversão de dependência (SOLID - DIP).
package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

// User é a entidade de negócio, independente de como é persistida.
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string // hash bcrypt, nunca a senha em texto plano
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository define as operações de persistência do usuário.
// A implementação concreta fica em internal/modules/user/repository.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*User, int, error)
}

// Service define as regras de negócio disponíveis para o módulo user.
// A implementação concreta fica em internal/modules/user/service.
type Service interface {
	Create(ctx context.Context, name, email, password string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, perPage int) ([]*User, int, error)
}
