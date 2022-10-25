package transporter

import (
	"bytes"
	"encoding/json"
)

type PayloadMessage[T any] struct {
	Data T `json:"data"`
}

func (pa *PayloadMessage[T]) DataBytes(data T) ([]byte, error) {
	var dataBytes bytes.Buffer

	err := json.NewEncoder(&dataBytes).Encode(data)
	if err != nil {
		return make([]byte, 0), nil
	}

	return dataBytes.Bytes(), nil
}

func NewPayloadMessage[T any](data T) PayloadMessage[T] {
	return PayloadMessage[T]{
		Data: data,
	}
}
