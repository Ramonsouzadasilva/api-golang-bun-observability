package handler

// TokenResponse é a resposta padrão contendo o par de tokens emitido.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// NewTokenResponse monta a resposta de tokens a partir do par emitido.
func NewTokenResponse(access, refresh string) TokenResponse {
	return TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}
}
