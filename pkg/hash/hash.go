package hash

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"go-api-base/internal/observability"
)

// Cost permite configurar o custo do bcrypt (menor para testes).
var Cost = bcrypt.DefaultCost

// Hash gera o hash bcrypt de uma senha em texto plano.
func Hash(password string) (string, error) {
	start := time.Now()
	defer func() {
		observability.BcryptDuration.Observe(time.Since(start).Seconds())
	}()

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare verifica se a senha em texto plano corresponde ao hash informado.
func Compare(hashedPassword, password string) bool {
	start := time.Now()
	defer func() {
		observability.BcryptDuration.Observe(time.Since(start).Seconds())
	}()

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
