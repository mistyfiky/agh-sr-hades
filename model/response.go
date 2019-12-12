package model

import "encoding/json"

type response struct {
	Meta *Meta `json:"meta"`
	Data data  `json:"data"`
}

func newResponseSuccess(data data, message *string) *response {
	return &response{Meta: newMetaSuccess(message), Data: data}
}

func NewResponseError(message string) *response {
	return &response{Meta: newMetaError(message)}
}

func (response *response) ToJson() []byte {
	result, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	return result
}
