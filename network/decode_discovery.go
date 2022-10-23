package network

import (
	"encoding/json"

	"github.com/block-api/block-node/transporter"
)

func decodeDiscoveryPayload(pocket *transporter.Pocket[[]byte]) (transporter.PayloadDiscovery, error) {
	disocveryPayload := transporter.PayloadDiscovery{}

	err := json.Unmarshal(pocket.Payload, &disocveryPayload)
	if err != nil {
		return disocveryPayload, err
	}

	return disocveryPayload, nil
}
