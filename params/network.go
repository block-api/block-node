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
	"fmt"

	"github.com/block-api/block-node/network/transport"
	"gopkg.in/yaml.v3"
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
	HeartbeatInterval int                 `yaml:"heartbeat_interval" json:"hi"`
	ActionTimeout     int                 `yaml:"action_timeout" json:"at"`
	Transport         transport.Type      `yaml:"transport" json:"t"`
	Settings          interface{}         `yaml:"settings" json:"s"`
	Nodes             []NetworkNodeConfig `yamls:"nodes" json:"n"`
}

type NetworkNodeConfig struct {
	NodeID     string         `yaml:"node_id" json:"nid"`
	Transport  transport.Type `yaml:"transport" json:"t"`
	PublicHost string         `yaml:"public_host" json:"ph"`
	PublicPort string         `yaml:"public_port" json:"pp"`
}

type NetworkTCPSettings struct {
	PublicHost string `yaml:"public_host" json:"ph"`
	PublicPort string `yaml:"public_port" json:"pp"`
	BindHost   string `yaml:"bind_host" json:"bh"`
	BindPort   string `yaml:"bind_port" json:"bp"`
}
type NetworkRedisSettings struct{}

func GetNetworkSettings[T NetworkTCPSettings | NetworkRedisSettings](config *NetworkConfig) (T, error) {
	var settings T

	fmt.Println(config.Settings)
	settingsBytes, err := yaml.Marshal(config.Settings)
	if err != nil {
		return settings, err
	}

	err = yaml.Unmarshal(settingsBytes, &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
