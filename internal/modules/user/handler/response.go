package handler

import (
	"time"

	"github.com/google/uuid"

	"go-api-base/internal/modules/user/domain"
	"go-api-base/pkg/pagination"
)

// UserResponse é a representação pública de um usuário (sem a senha).
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse converte a entidade de domínio na resposta pública da API.
func NewUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserListResponse é a resposta paginada de listagem de usuários.
type UserListResponse struct {
	Data []UserResponse  `json:"data"`
	Meta pagination.Meta `json:"meta"`
}

// NewUserListResponse monta a resposta paginada a partir da lista de usuários.
func NewUserListResponse(users []*domain.User, page, perPage, total int) UserListResponse {
	data := make([]UserResponse, len(users))
	for i, u := range users {
		data[i] = NewUserResponse(u)
	}
	return UserListResponse{Data: data, Meta: pagination.NewMeta(page, perPage, total)}
}
