package network

import (
	"bytes"
	"encoding/json"

	"github.com/block-api/block-node/transporter"
	"github.com/ethereum/go-ethereum/crypto"
)

type Pocket struct {
	channel transporter.Channel
	from    string
	target  string
	payload interface{}
	hash    []byte
}

func (p *Pocket) Channel() transporter.Channel {
	return p.channel
}

func (p *Pocket) From() string {
	return p.from
}

func (p *Pocket) Payload() interface{} {
	return p.payload
}

type Event uint

const (
	EventConnected Event = iota + 1
	EventDisconnected
	EventPing
	EventPong
)

type PayloadDiscovery struct {
	NodeID  string              `json:"node_id"`
	Name    string              `json:"name"`
	Blocks  map[string][]string `json:"blocks"`
	Version uint                `json:"version"`
	Event   Event               `json:"event"`
}

type PayloadMessage interface{}

// NewPocket creates new network pocket
func NewPocket[P PayloadDiscovery | PayloadMessage](channel transporter.Channel, from string, target string, payload P) Pocket {
	var payloadBytes bytes.Buffer

	json.NewEncoder(&payloadBytes).Encode(payload)
	hash := crypto.Keccak256(payloadBytes.Bytes())

	return Pocket{
		channel: channel,
		from:    from,
		target:  target,
		payload: payload,
		hash:    hash,
	}
}
