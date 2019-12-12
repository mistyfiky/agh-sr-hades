package model

type SimpleResponse struct {
	*response
}

func NewSimpleResponse(message string) *SimpleResponse {
	return &SimpleResponse{response: newResponseSuccess(nil, &message)}
}
