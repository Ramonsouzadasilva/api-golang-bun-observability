package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registra as rotas do módulo user no router informado.
// Todas as rotas exigem autenticação (authMW).
func (m *Module) RegisterRoutes(r chi.Router, authMW func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Use(authMW)

		r.Get("/", m.Handler.List)
		r.Get("/me", m.Handler.Me)
		r.Get("/{id}", m.Handler.GetByID)
		r.Put("/{id}", m.Handler.Update)
		r.Delete("/{id}", m.Handler.Delete)
	})
}
