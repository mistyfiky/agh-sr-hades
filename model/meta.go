package model

type Meta struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}

func newMetaSuccess(message *string) *Meta {
	return &Meta{Success: true, Message: message}
}

func newMetaError(message string) *Meta {
	return &Meta{Message: &message}
}
