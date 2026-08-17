package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"go-api-base/internal/modules/auth/domain"
)

// RefreshToken é o modelo de persistência da tabela "refresh_tokens".
type RefreshToken struct {
	bun.BaseModel `bun:"table:refresh_tokens,alias:rt"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserID    uuid.UUID  `bun:"user_id,notnull,type:uuid"`
	Token     string     `bun:"token,notnull,unique"`
	ExpiresAt time.Time  `bun:"expires_at,notnull"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	RevokedAt *time.Time `bun:"revoked_at"`
}

// ToDomain converte o modelo de persistência para a entidade de domínio.
func (m *RefreshToken) ToDomain() *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
		RevokedAt: m.RevokedAt,
	}
}

// FromDomain converte a entidade de domínio para o modelo de persistência.
func FromDomain(rt *domain.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
	}
}
