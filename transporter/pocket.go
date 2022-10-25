package transporter

import (
	"bytes"
	"encoding/json"

	"github.com/block-api/block-node/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

type Pocket[P []byte | any] struct {
	Channel      Channel      `json:"channel"`
	FromID       utils.NodeID `json:"from_id"`
	TargetID     utils.NodeID `json:"target_id"`
	TargetAction string       `json:"target_action"`
	Payload      P            `json:"payload"`
	Hash         string       `json:"hash"`
}

// NewPocket creates new network pocket
func NewPocket[P PayloadDiscovery | PayloadMessage[any]](channel Channel, fromID utils.NodeID, targetID utils.NodeID, payload P) Pocket[[]byte] {
	var payloadBytes bytes.Buffer
	pocket := Pocket[[]byte]{
		Channel:  channel,
		FromID:   fromID,
		TargetID: targetID,
	}

	json.NewEncoder(&payloadBytes).Encode(payload)
	hash := crypto.Keccak256(payloadBytes.Bytes())

	pocket.Payload = payloadBytes.Bytes()
	pocket.Hash = crypto.Keccak256Hash(hash).String()

	return pocket
}
