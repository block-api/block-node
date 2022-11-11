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
	"fmt"
	"sync"

	"github.com/block-api/block-node/log"
)

func (m *Manager) senderWorker(cStop <-chan int, cSend <-chan Packet, wgStop *sync.WaitGroup) {
	log.Debug("network::sender_worker::start")
L:
	for {
		select {
		case <-cStop:
			log.Debug("network::sender_worker::stop")
			wgStop.Done()
			break L
		case packet := <-cSend:
			log.Debug("network::sender_worker::received_packet")
			fmt.Println(packet)
			fmt.Println("-- send packet --")

			go sendPacketWorker(packet)

			continue
		}
	}
}

func sendPacketWorker(packet Packet) {
	log.Debug("network::send_packet_worker::start")
}
