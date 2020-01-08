package model

type User struct {
	Username string `json:"username"`
	Movies []string `json:"movies"`
}

type UserResponse struct {
	Meta Meta `json:"meta"`
	Data User `json:"data"`
}
