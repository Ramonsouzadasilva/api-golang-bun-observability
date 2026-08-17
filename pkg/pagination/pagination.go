// Package pagination centraliza as regras de paginação usadas nas
// listagens de todos os módulos da API.
package pagination

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Meta é o bloco de metadados de paginação retornado nas respostas de listagem.
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Normalize aplica os valores padrão e limites de página/itens por página.
func Normalize(page, perPage int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}
	return page, perPage
}

// LimitOffset converte page/perPage em limit/offset para uso nas queries SQL.
func LimitOffset(page, perPage int) (limit, offset int) {
	page, perPage = Normalize(page, perPage)
	return perPage, (page - 1) * perPage
}

// NewMeta monta os metadados de paginação a partir do total de registros.
func NewMeta(page, perPage, total int) Meta {
	page, perPage = Normalize(page, perPage)

	totalPages := total / perPage
	if total%perPage != 0 {
		totalPages++
	}

	return Meta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}
