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

// Package packet
package packet

import "github.com/block-api/block-node/network/transport"

type Node struct {
	Transport  transport.Type  `yaml:"transport" json:"t"`
	NodeID     string          `yaml:"node_id" json:"nid"`
	PublicHost string          `yaml:"public_host" json:"ph"`
	PublicPort string          `yaml:"public_port" json:"pp"`
	Functions  map[string]bool `yaml:"functions" json:"f"`
}
