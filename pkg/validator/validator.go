// Package validator centraliza a validação de DTOs de entrada usando
// tags `validate` (go-playground/validator), com mensagens de erro
// legíveis em vez do erro cru da lib.
package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate valida a struct informada de acordo com as tags `validate`
// dos seus campos, retornando um erro legível com todas as violações.
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	messages := make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		messages = append(messages, formatFieldError(e))
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// DecodeAndValidate decodifica o corpo JSON da requisição para dst e,
// em seguida, valida os campos preenchidos.
func DecodeAndValidate(r *http.Request, dst interface{}) error {
	// Limita o tamanho do body a 1MB para evitar ataques de exaustão de memória
	body := io.LimitReader(r.Body, 1048576)

	if err := json.NewDecoder(body).Decode(dst); err != nil {
		return fmt.Errorf("invalid request body")
	}
	return Validate(dst)
}

func formatFieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
