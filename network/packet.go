// Copyright 2022 The block-node Authors
// This file is part of the block-node library.
//
// The block-node library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The block-node library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the block-node library. If not, see <http://www.gnu.org/licenses/>.

// Package network
package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/block-api/block-node/network/delivery"
	"github.com/block-api/block-node/network/packet"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrInvalidHeader = []byte("INVALID_HEADER")
var Header = []byte(`\xblock-node\xv0.2.0`)

type Packet struct {
	Delivery           delivery.Type `json:"d"`
	Hash               []byte        `json:"h"`
	ResponseHash       []byte        `json:"rh,omitempty"`
	Type               packet.Type   `json:"t"`
	FromID             string        `json:"fid"`
	TargetID           string        `json:"tid,omitempty"`
	TargetNodeFunction string        `json:"tnf,omitempty"`
	Body               any           `json:"b,omitempty"`
	CreatedAt          int64         `json:"cat"`
}

func NewPacket(deliveryType delivery.Type, packetType packet.Type, fromID string, targetID string, targetFunction string, body interface{}, responseHash *[]byte) Packet {
	var respHash []byte
	if responseHash != nil {
		respHash = *responseHash
	}

	newPacket := Packet{
		ResponseHash:       respHash,
		Delivery:           deliveryType,
		Type:               packetType,
		FromID:             fromID,
		TargetID:           targetID,
		TargetNodeFunction: targetFunction,
		Body:               body,
		CreatedAt:          time.Now().UnixMilli(),
	}

	bodyBytes, _ := json.Marshal(body)
	hash := crypto.Keccak256(bodyBytes, []byte(newPacket.Delivery), []byte(newPacket.Type), []byte(newPacket.FromID), []byte(newPacket.TargetID), []byte(strconv.FormatInt(newPacket.CreatedAt, 10)))
	newPacket.Hash = hash

	return newPacket
}

func (p *Packet) Bytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return make([]byte, 0), err
	}

	return append(Header, jsonBytes...), nil
}

func (p *Packet) Validate() error {
	bodyBytes, _ := json.Marshal(p.Body)

	hash := crypto.Keccak256(bodyBytes, []byte(p.Delivery), []byte(p.Type), []byte(p.FromID), []byte(p.TargetID), []byte(strconv.FormatInt(p.CreatedAt, 10)))

	if bytes.Compare(hash, p.Hash) == 0 {
		return nil
	}

	return errors.New("invalid hash")
}

func EncodePacket(data interface{}) []byte {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return make([]byte, 0)
	}

	return append(Header, jsonBytes...)
}

func DecodePacket(data []byte) Packet {
	var decodedPacket Packet
	_ = json.Unmarshal(data, &decodedPacket)

	switch decodedPacket.Type {
	case packet.Heartbeat:
		var bodyHeartbeat packet.HeartbeatBody

		bodyBytes, _ := json.Marshal(decodedPacket.Body)
		_ = json.Unmarshal(bodyBytes, &bodyHeartbeat)
		decodedPacket.Body = bodyHeartbeat

		break
	case packet.Cmd:
		var bodyCmd packet.CmdBody

		bodyBytes, _ := json.Marshal(decodedPacket.Body)
		_ = json.Unmarshal(bodyBytes, &bodyCmd)
		decodedPacket.Body = bodyCmd

		break
	case packet.Function:
		var bodyCmd packet.FunctionBody

		bodyBytes, _ := json.Marshal(decodedPacket.Body)
		_ = json.Unmarshal(bodyBytes, &bodyCmd)
		decodedPacket.Body = bodyCmd

		break
	default:
		break
	}

	return decodedPacket
}
