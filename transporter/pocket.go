package transporter

import (
	"bytes"
	"encoding/json"

	"github.com/block-api/block-node/common/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Pocket[P []byte | any] struct {
	Channel      Channel               `json:"channel"`
	VersionName  types.NodeVersionName `json:"version_name"`
	FromID       types.NodeID          `json:"from_id"`
	TargetID     *types.NodeID         `json:"target_id"`
	TargetAction *types.TargetAction   `json:"target_action"`
	Payload      P                     `json:"payload"`
	Hash         string                `json:"hash"`
}

// NewPocket creates new network pocket
func NewPocket[P PayloadDiscovery | PayloadMessage](channel Channel, versionName types.NodeVersionName, fromID types.NodeID, targetID *types.NodeID, targetAction *types.TargetAction, payload P) Pocket[[]byte] {
	var payloadBytes bytes.Buffer
	pocket := Pocket[[]byte]{
		Channel:      channel,
		VersionName:  versionName,
		FromID:       fromID,
		TargetID:     targetID,
		TargetAction: targetAction,
	}

	json.NewEncoder(&payloadBytes).Encode(payload)
	hash := crypto.Keccak256(payloadBytes.Bytes())

	pocket.Payload = payloadBytes.Bytes()
	pocket.Hash = crypto.Keccak256Hash(hash).String()

	return pocket
}
