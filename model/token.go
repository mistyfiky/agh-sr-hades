package model

type Token struct {
	Token []byte `json:"token"`
}

func newToken(token []byte) *Token {
	return &Token{Token: token}
}

type TokenResponse struct {
	*response
	Data *Token `json:"data"`
}

func NewTokenResponse(token []byte) *TokenResponse {
	return &TokenResponse{response: newResponseSuccess(newToken(token), nil)}
}
