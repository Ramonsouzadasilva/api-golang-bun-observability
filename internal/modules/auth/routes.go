package auth

import "github.com/go-chi/chi/v5"

// RegisterRoutes registra as rotas públicas do módulo auth no router informado.
func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", m.Handler.Register)
		r.Post("/login", m.Handler.Login)
		r.Post("/refresh", m.Handler.Refresh)
		r.Post("/logout", m.Handler.Logout)
	})
}
