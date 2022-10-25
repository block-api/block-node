package network

import (
	"encoding/json"

	"github.com/block-api/block-node/transporter"
)

func decodePayload[T transporter.PayloadDiscovery | transporter.PayloadMessage[any]](pocket *transporter.Pocket[[]byte]) (T, error) {
	var payload T

	err := json.Unmarshal(pocket.Payload, &payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}
