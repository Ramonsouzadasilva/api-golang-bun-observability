// Package user é o ponto de entrada público do módulo: monta (wire) o
// repository, o service e o handler, e expõe as rotas HTTP do módulo.
package user

import (
	"github.com/uptrace/bun"

	"go-api-base/internal/modules/user/domain"
	"go-api-base/internal/modules/user/handler"
	"go-api-base/internal/modules/user/repository"
	"go-api-base/internal/modules/user/service"
)

// Module agrupa as dependências construídas do módulo user, prontas
// para serem usadas pelo router e por outros módulos (ex: auth).
type Module struct {
	Handler    *handler.UserHandler
	Service    domain.Service
	Repository domain.Repository
}

// NewModule constrói o módulo user por completo a partir da conexão com o banco.
func NewModule(db *bun.DB) *Module {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	return &Module{
		Handler:    h,
		Service:    svc,
		Repository: repo,
	}
}
