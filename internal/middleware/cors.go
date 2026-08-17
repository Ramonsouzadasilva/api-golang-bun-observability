package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CORS retorna um middleware de CORS com uma configuração permissiva,
// adequada para uma API base. Ajuste AllowedOrigins em produção.
func CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
