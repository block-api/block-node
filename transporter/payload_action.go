package transporter

import (
	"bytes"
	"encoding/json"
)

type PayloadAction[T any] struct {
	Data T `json:"data"`
}

func (pa *PayloadAction[T]) DataBytes(data T) ([]byte, error) {
	var dataBytes bytes.Buffer

	err := json.NewEncoder(&dataBytes).Encode(data)
	if err != nil {
		return make([]byte, 0), nil
	}

	return dataBytes.Bytes(), nil
}

func NewPayloadAction[T any](data T) PayloadAction[T] {
	return PayloadAction[T]{
		Data: data,
	}
}
