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
package network

import (
	"sync"

	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network/delivery"
	"github.com/block-api/block-node/network/packet"
	"github.com/block-api/block-node/network/transport"
	"github.com/block-api/block-node/params"
)

func (m *Manager) networkWorker(cStop <-chan int, wgStop *sync.WaitGroup) {
	log.Debug("network::network_worker::start")
L:
	for {
		select {
		case <-cStop:
			break L
		}
	}

	log.Debug("network::network_worker::stop")
	wgStop.Done()
}

type Heartbeat struct {
	Packet
}

func NewHeartbeat(config *params.NetworkConfig, fromID string, targetID string, targetFunction string) Packet {
	knownNodes := make(map[string]packet.Node)

	for nodeID, node := range manager.router.knownNodes {
		knownNodes[nodeID] = packet.Node{
			NodeID:     node.NodeID,
			Transport:  node.Transport,
			PublicHost: node.PublicHost,
			PublicPort: node.PublicPort,
		}
	}

	beatBody := packet.HeartbeatBody{
		NodeID:     fromID,
		Transport:  config.Transport,
		KnownNodes: knownNodes,
	}

	if config.Transport == transport.TCP {
		tcpSettings, _ := params.GetNetworkSettings[params.NetworkTCPSettings](config)
		beatBody.PublicHost = tcpSettings.PublicHost
		beatBody.PublicPort = tcpSettings.PublicPort
	}

	return NewPacket(delivery.All, packet.Heartbeat, fromID, targetID, targetFunction, beatBody, nil)
}
