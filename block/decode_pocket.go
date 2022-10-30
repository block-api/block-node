package block

import (
	"encoding/json"

	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

func DecodePocket(pocketBytes []byte) (transporter.Pocket[[]byte], error) {
	pocket := transporter.Pocket[[]byte]{}

	err := json.Unmarshal(pocketBytes, &pocket)
	if err != nil {
		log.Warning(err.Error())
		return transporter.Pocket[[]byte]{}, err
	}

	return pocket, nil
}

func DecodePayload[T interface{}](payloadBytes []byte) (T, error) {
	var payload T

	err := json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}
