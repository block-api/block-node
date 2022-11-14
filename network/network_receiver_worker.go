// Package network
//
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
	"fmt"
	"sync"

	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network/packet"
)

func (m *Manager) receiverWorker(transport Transport, cStop <-chan int, cReceive chan Packet, wgStop *sync.WaitGroup) {
	log.Debug("network::receiver_worker::start")

	go m.transport.Start(cReceive)
L:
	for {
		select {
		case <-cStop:
			log.Debug("network::receiver_worker::stop")
			// wgStop.Done()
			break L
		case receivedPacket := <-cReceive:
			fmt.Println(receivedPacket)

			log.Debug("network::receiver_worker::received_packet")
			fmt.Println(receivedPacket)
			fmt.Println("^^^^^^^^^^^^^^")

			if receivedPacket.Type == packet.Heartbeat {
				log.Debug("-- process heartbeat --")

				return
			}

			fmt.Println("----")
			continue
		}
	}
}
