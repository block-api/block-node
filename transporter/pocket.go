package transporter

import (
	"bytes"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
)

type Pocket[P []byte | interface{}] struct {
	Channel Channel `json:"channel"`
	From    string  `json:"from"`
	Target  string  `json:"target"`
	Payload P       `json:"payload"`
	Hash    string  `json:"hash"`
}

// NewPocket creates new network pocket
func NewPocket[P PayloadDiscovery | PayloadMessage](channel Channel, from string, target string, payload P) Pocket[[]byte] {
	var payloadBytes bytes.Buffer
	pocket := Pocket[[]byte]{
		Channel: channel,
		From:    from,
		Target:  target,
	}

	json.NewEncoder(&payloadBytes).Encode(payload)
	hash := crypto.Keccak256(payloadBytes.Bytes())

	pocket.Payload = payloadBytes.Bytes()
	pocket.Hash = crypto.Keccak256Hash(hash).String()

	return pocket
}
