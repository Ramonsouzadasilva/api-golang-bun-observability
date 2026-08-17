package handler

// UpdateUserRequest é o payload aceito para atualização de um usuário.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}
