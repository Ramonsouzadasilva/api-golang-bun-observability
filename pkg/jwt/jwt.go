// Package jwt encapsula a geração e validação de access tokens JWT.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken é retornado quando o token é inválido, malformado ou expirado.
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims representa as informações carregadas dentro do token.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// Manager gera e valida tokens JWT assinados com HMAC.
type Manager struct {
	secret []byte
}

// NewManager cria um Manager a partir do segredo configurado na aplicação.
func NewManager(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

// Generate cria um novo access token para o usuário informado, válido pelo
// tempo (ttl) especificado.
func (m *Manager) Generate(userID uuid.UUID, email string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse valida o token e retorna as claims nele contidas.
func (m *Manager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
