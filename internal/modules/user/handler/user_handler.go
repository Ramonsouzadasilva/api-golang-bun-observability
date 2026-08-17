// Package handler contém a camada HTTP do módulo user: decodifica
// requests, chama o service e escreve as responses. Não conhece nada
// sobre banco de dados nem SQL.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"go-api-base/internal/middleware"
	"go-api-base/internal/modules/user/domain"
	"go-api-base/pkg/httpresponse"
	"go-api-base/pkg/validator"
)

// UserHandler expõe os endpoints HTTP do módulo user.
type UserHandler struct {
	service domain.Service
}

// NewUserHandler cria um UserHandler a partir do service do módulo.
func NewUserHandler(service domain.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Me retorna os dados do usuário autenticado (extraído do token JWT).
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		httpresponse.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, NewUserResponse(u))
}

// GetByID retorna um usuário pelo ID.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, NewUserResponse(u))
}

// Update atualiza os dados de um usuário.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		httpresponse.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.service.Update(r.Context(), id, req.Name)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, NewUserResponse(u))
}

// Delete remove (soft delete) um usuário.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusNoContent, nil)
}

// List retorna uma listagem paginada de usuários.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	users, total, err := h.service.List(r.Context(), page, perPage)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, NewUserListResponse(users, page, perPage, total))
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		httpresponse.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrEmailAlreadyInUse):
		httpresponse.Error(w, http.StatusConflict, err.Error())
	default:
		httpresponse.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
