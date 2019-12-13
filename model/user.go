package model

type User struct {
	Username string `json:"username"`
}

type UserResponse struct {
	Meta Meta `json:"meta"`
	Data User `json:"data"`
}
