// Package router monta o http.Handler final da aplicação: middlewares
// globais, rota de health check e as rotas de cada módulo.
package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"go-api-base/internal/middleware"
	authmodule "go-api-base/internal/modules/auth"
	usermodule "go-api-base/internal/modules/user"
	"go-api-base/pkg/httpresponse"
	"go-api-base/pkg/jwt"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Dependencies contém tudo que o router precisa para montar as rotas.
type Dependencies struct {
	Logger     *slog.Logger
	JWTManager *jwt.Manager
	AuthModule *authmodule.Module
	UserModule *usermodule.Module
}

func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(middleware.RequestLogger(deps.Logger))
	r.Use(middleware.CORS())
	r.Use(middleware.Metrics())
	r.Use(chimw.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpresponse.JSON(
			w,
			http.StatusOK,
			map[string]string{"status": "ok"},
		)
	})

	r.Handle("/metrics", promhttp.Handler())

	authMW := middleware.Auth(deps.JWTManager)

	r.Route("/api/v1", func(r chi.Router) {
		deps.AuthModule.RegisterRoutes(r)
		deps.UserModule.RegisterRoutes(r, authMW)
	})

	return r
}
