package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-api-base/pkg/httpresponse"
	"go-api-base/pkg/jwt"
)

type contextKey string

const (
	// UserIDKey é a chave de contexto onde o ID do usuário autenticado é armazenado.
	UserIDKey contextKey = "user_id"
	// UserEmailKey é a chave de contexto onde o e-mail do usuário autenticado é armazenado.
	UserEmailKey contextKey = "user_email"
)

// Auth retorna um middleware que exige um Bearer token JWT válido no
// header Authorization, expondo o usuário autenticado no contexto da requisição.
func Auth(manager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httpresponse.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				httpresponse.Error(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			claims, err := manager.Parse(parts[1])
			if err != nil {
				httpresponse.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
