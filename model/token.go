package model

type Token struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Meta Meta  `json:"meta"`
	Data Token `json:"data"`
}
