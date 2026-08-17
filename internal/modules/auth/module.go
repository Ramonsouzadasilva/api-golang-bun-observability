// Package auth é o ponto de entrada público do módulo de autenticação.
// Depende da interface userdomain.Repository (não da implementação
// concreta do módulo user), mantendo os módulos desacoplados.
package auth

import (
	"github.com/uptrace/bun"

	"go-api-base/internal/config"
	"go-api-base/internal/modules/auth/domain"
	"go-api-base/internal/modules/auth/handler"
	"go-api-base/internal/modules/auth/repository"
	"go-api-base/internal/modules/auth/service"
	userdomain "go-api-base/internal/modules/user/domain"
	"go-api-base/pkg/jwt"
)

// Module agrupa as dependências construídas do módulo auth.
type Module struct {
	Handler *handler.AuthHandler
	Service domain.Service
}

// NewModule constrói o módulo auth. Recebe o repository de usuários do
// módulo user para não precisar reimplementar acesso a dados de usuário.
func NewModule(
	db *bun.DB,
	userRepo userdomain.Repository,
	jwtManager *jwt.Manager,
	jwtCfg config.JWTConfig,
) *Module {
	repo := repository.NewAuthRepository(db)
	svc := service.NewAuthService(repo, userRepo, jwtManager, jwtCfg.AccessTokenTTL, jwtCfg.RefreshTokenTTL)
	h := handler.NewAuthHandler(svc)

	return &Module{
		Handler: h,
		Service: svc,
	}
}
