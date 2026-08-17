# go-api-base

API base em Go para servir como ponto de partida de novas APIs. Já vem com
módulos de **user** e **auth** (JWT + refresh token), estrutura modular,
Docker Compose com Postgres e as ferramentas mais usadas do ecossistema Go
para APIs REST.

## Stack

- **[chi](https://github.com/go-chi/chi)** — router HTTP
- **[pgx](https://github.com/jackc/pgx)** — driver Postgres (usado como driver `database/sql`)
- **[bun](https://github.com/uptrace/bun)** — query builder / ORM, sobre o pgx
- **[uuid](https://github.com/google/uuid)** — identificadores das entidades
- **[golang-jwt](https://github.com/golang-jwt/jwt)** — geração/validação de access tokens
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** — hash de senhas
- **[validator](https://github.com/go-playground/validator)** — validação de DTOs
- **[log/slog](https://pkg.go.dev/log/slog)** (stdlib) — logs estruturados
- **Docker Compose** — API + Postgres

## Arquitetura

Cada módulo de negócio (`user`, `auth`, ...) segue a mesma organização em
camadas, aplicando os princípios SOLID — principalmente **inversão de
dependência**: `service` e `handler` dependem apenas das interfaces
definidas em `domain`, nunca de implementações concretas.

```
internal/modules/<modulo>/
    domain/       -> entidade de negócio + interfaces Repository e Service (contratos)
    model/        -> struct mapeado para o banco via bun (tags de coluna)
    repository/   -> implementação de domain.Repository (SQL fica só aqui)
    service/      -> implementação de domain.Service (regras de negócio)
    handler/      -> camada HTTP: request.go, response.go e os handlers
    routes.go     -> registra as rotas do módulo no chi.Router
    module.go     -> "wiring": injeta repository -> service -> handler
```

> **Nota de organização:** no pedido original, `request.go`/`response.go`
> ficavam na raiz do módulo. Eles foram movidos para dentro de `handler/`
> porque `module.go` (raiz do módulo) importa o pacote `handler`, e se
> `handler` importasse de volta a raiz para pegar os DTOs, teríamos um
> import cycle (proibido em Go). Colocando os DTOs dentro de `handler/`,
> a dependência fica só em uma direção. `routes.go` e `module.go`
> continuam na raiz do módulo, como pedido.

Módulos podem depender de **interfaces** de outros módulos sem acoplamento
direto: o módulo `auth` recebe o `userdomain.Repository` do módulo `user`
via injeção de dependência em `NewModule`, sem importar a implementação
concreta do repositório de usuários.

### Autenticação

- **Access token**: JWT (HS256), curta duração (padrão 15 min), enviado no
  header `Authorization: Bearer <token>`.
- **Refresh token**: string opaca aleatória, persistida na tabela
  `refresh_tokens`, com expiração longa (padrão 7 dias) e suporte a
  revogação (logout) e rotação (a cada `/auth/refresh`, o token antigo é
  revogado e um novo é emitido).

### Soft delete

O model de `User` usa a tag `bun:"...,soft_delete"` no campo `DeletedAt`.
Isso faz o bun tratar automaticamente:
- `DELETE` → vira `UPDATE ... SET deleted_at = now()`
- `SELECT` → ignora automaticamente registros com `deleted_at` preenchido

## Endpoints

| Método | Rota                  | Autenticado | Descrição                          |
|--------|-----------------------|:-----------:|-------------------------------------|
| GET    | `/health`              | não         | Health check                        |
| POST   | `/api/v1/auth/register`| não         | Cria conta e retorna tokens         |
| POST   | `/api/v1/auth/login`   | não         | Login, retorna tokens               |
| POST   | `/api/v1/auth/refresh` | não         | Renova o access token               |
| POST   | `/api/v1/auth/logout`  | não         | Revoga o refresh token              |
| GET    | `/api/v1/users`        | sim         | Lista usuários (paginado)           |
| GET    | `/api/v1/users/me`     | sim         | Dados do usuário autenticado        |
| GET    | `/api/v1/users/{id}`   | sim         | Busca usuário por ID                |
| PUT    | `/api/v1/users/{id}`   | sim         | Atualiza nome do usuário            |
| DELETE | `/api/v1/users/{id}`   | sim         | Remove (soft delete) usuário        |

## Como rodar

### 1. Com Docker (recomendado)

```bash
cp .env.example .env
make docker-up          # sobe API + Postgres
make migrate-up         # roda as migrations
make seed               # (opcional) cria o usuário admin@example.com / Admin@123
```

A API sobe em `http://localhost:8080`.

### 2. Localmente (sem Docker)

Requer Go 1.22+ e um Postgres rodando.

```bash
cp .env.example .env    # ajuste DB_HOST etc se necessário
go mod tidy              # baixa as dependências (necessário: sem acesso à
                          # internet neste ambiente, o go.sum não foi gerado)
make migrate-up
make seed                 # opcional
make run
```

### Exemplo de uso

```bash
# registrar
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"supersecret"}'

# login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"supersecret"}'

# usar o access_token retornado
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <access_token>"
```

## Comandos úteis (Makefile)

| Comando               | Descrição                                   |
|------------------------|----------------------------------------------|
| `make run`             | Roda a API localmente                        |
| `make build`           | Compila o binário em `bin/api`               |
| `make test`            | Roda os testes com cobertura                 |
| `make tidy`            | `go mod tidy`                                |
| `make migrate-up`      | Aplica migrations pendentes                  |
| `make migrate-down`    | Desfaz o último grupo de migrations          |
| `make migrate-status`  | Mostra status das migrations                 |
| `make seed`            | Popula o banco com dados iniciais            |
| `make docker-up`       | Sobe API + Postgres via Docker Compose       |
| `make docker-down`     | Derruba os containers                        |
| `make docker-logs`     | Acompanha logs da API                        |

## Adicionando um novo módulo

1. Copie a estrutura de `internal/modules/user` como referência.
2. Defina a entidade e as interfaces em `domain/`.
3. Implemente `model/`, `repository/` e `service/`.
4. Implemente `handler/` (request.go, response.go, handlers) usando apenas
   `domain.Service`.
5. Crie `routes.go` e `module.go` na raiz do módulo.
6. Registre o módulo em `internal/router/router.go` e em `cmd/api/main.go`.
7. Adicione as migrations SQL correspondentes em
   `internal/database/migrations/`.

## Sobre o go.sum

Este projeto foi gerado em um ambiente sem acesso à internet, então o
`go.sum` não pôde ser gerado (ele depende de baixar os módulos para
calcular os hashes). Rode `go mod tidy` na primeira vez que abrir o
projeto — isso vai baixar as dependências do `go.mod` e gerar o `go.sum`
automaticamente. O `Dockerfile` já espera um `go.sum` presente antes do
build (`COPY go.mod go.sum ./`), então rode `go mod tidy` antes do
primeiro `make docker-up`.
