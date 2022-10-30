package transporter

import (
	"bytes"
	"encoding/json"
)

type PayloadMessage struct {
	Data any `json:"data"`
}

func (pa *PayloadMessage) DataBytes(data any) ([]byte, error) {
	var dataBytes bytes.Buffer

	err := json.NewEncoder(&dataBytes).Encode(data)
	if err != nil {
		return make([]byte, 0), nil
	}

	return dataBytes.Bytes(), nil
}

func (pa *PayloadMessage) JSON() ([]byte, error) {
	json, err := json.Marshal(pa)
	if err != nil {
		return make([]byte, 0), err
	}

	return json, nil
}

func NewPayloadMessage(data any) PayloadMessage {
	return PayloadMessage{
		Data: data,
	}
}
