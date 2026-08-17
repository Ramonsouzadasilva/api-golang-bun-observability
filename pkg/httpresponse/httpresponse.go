// Package httpresponse padroniza a forma como a API escreve respostas JSON,
// evitando repetição de encoding/headers em cada handler.
package httpresponse

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse é o formato padrão usado para retornar erros da API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// JSON escreve o payload informado como JSON, com o status HTTP indicado.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// Error escreve uma resposta de erro padronizada.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}
