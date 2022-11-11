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
package params

import (
	"encoding/json"

	"github.com/block-api/block-node/network/transport"
)

// DefaultNetworkTransport default transport method used by block node network module
const DefaultNetworkTransport = transport.TCP

// Default bind host
const DefaultNetworkBindHost = "0.0.0.0"

// Default bind port
const DefaultNetworkBindPort = "8090"

// DefaultActionTimeout defines default time out for block actions.
// If no timeout was defined in config file this value will be used
const DefaultNetworkActionTimeout int = 15

// DefaultHeartbeatInterval defines how often node will be sending
// heartbeat packet to other known nodes in the network. Nodes which
// are not sending heartbeat peridiocaly will be automatically removed
// from other nodes network topology - hence will be not receiving requests
const DefaultNetworkHeartbeatInterval int = 5

// NetworkConfig describs configuration options for network
type NetworkConfig struct {
	HeartbeatInterval int            `yaml:"heartbeat_interval"`
	ActionTimeout     int            `yaml:"action_timeout"`
	Transport         transport.Type `yaml:"transport"`
	Settings          any            `yaml:"settings"`
}

type NetworkTCPSettings struct {
	PublicHost string `yaml:"public_host"`
	PublicPort string `yaml:"public_port"`
	BindHost   string `yaml:"bind_host"`
	BindPort   string `yaml:"bind_port"`
}
type NetworkRedisSettings struct{}

func GetNetworkSettings[T NetworkTCPSettings | NetworkRedisSettings](config *NetworkConfig) (*T, error) {
	settingsBytes, err := json.Marshal(config.Settings)
	if err != nil {
		return nil, err
	}

	var settings T

	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
