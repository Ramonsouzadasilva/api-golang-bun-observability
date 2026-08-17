// Package model contém os structs mapeados para o banco via bun.
// Ficam isolados do domain para que a entidade de negócio não precise
// conhecer detalhes de persistência (tags bun, nomes de coluna etc).
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"go-api-base/internal/modules/user/domain"
)

// User é o modelo de persistência da tabela "users".
// O campo DeletedAt com a tag soft_delete faz o bun tratar deleções
// como soft delete automaticamente (UPDATE em vez de DELETE físico).
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name      string     `bun:"name,notnull"`
	Email     string     `bun:"email,notnull,unique"`
	Password  string     `bun:"password,notnull"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero"`
}

// ToDomain converte o modelo de persistência para a entidade de negócio.
func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromDomain converte a entidade de negócio para o modelo de persistência.
func FromDomain(u *domain.User) *User {
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
