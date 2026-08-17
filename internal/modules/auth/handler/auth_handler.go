// Package handler contém a camada HTTP do módulo auth.
package handler

import (
	"errors"
	"net/http"

	"go-api-base/internal/modules/auth/domain"
	userdomain "go-api-base/internal/modules/user/domain"
	"go-api-base/pkg/httpresponse"
	"go-api-base/pkg/validator"
)

// AuthHandler expõe os endpoints HTTP de autenticação.
type AuthHandler struct {
	service domain.Service
}

// NewAuthHandler cria um AuthHandler a partir do service do módulo.
func NewAuthHandler(service domain.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register cria uma nova conta e já retorna o par de tokens autenticado.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		httpresponse.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusCreated, NewTokenResponse(tokens.AccessToken, tokens.RefreshToken))
}

// Login autentica um usuário existente e retorna o par de tokens.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		httpresponse.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, NewTokenResponse(tokens.AccessToken, tokens.RefreshToken))
}

// Refresh troca um refresh token válido por um novo par de tokens.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		httpresponse.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, NewTokenResponse(tokens.AccessToken, tokens.RefreshToken))
}

// Logout revoga um refresh token, encerrando a sessão.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		httpresponse.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		h.handleError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		httpresponse.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, userdomain.ErrEmailAlreadyInUse):
		httpresponse.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrRefreshTokenNotFound), errors.Is(err, domain.ErrRefreshTokenExpired):
		httpresponse.Error(w, http.StatusUnauthorized, err.Error())
	default:
		httpresponse.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
