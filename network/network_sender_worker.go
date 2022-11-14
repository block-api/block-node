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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/block-api/block-node/block/function"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/params"
)

var (
	ErrPacketRequestTimeout = errors.New("packet request timeout")
)

func (m *Manager) senderWorker(nodeID string, config *params.NetworkConfig, transpo *Transport, cStop <-chan int, cSend <-chan Packet, wgStop *sync.WaitGroup) {
	log.Debug("network::sender_worker::start")
L:
	for {
		select {
		case <-cStop:
			log.Debug("network::sender_worker::stop")
			// wgStop.Done()
			break L
		case packet := <-cSend:
			// todo: limiting rate X / s
			log.Debug("network::sender_worker::received_packet")
			go sendPacketWorker(m, config, packet, transpo)

			continue
		}
	}
}

func sendPacketWorker(m *Manager, config *params.NetworkConfig, sendPacket Packet, transpo *Transport) {
	log.Debug("network::send_packet_worker::start")
	var tickerTimeout = time.NewTicker(time.Second * time.Duration(config.ActionTimeout))

	if m.nodeID == sendPacket.TargetID && sendPacket.TargetNodeFunction != "" {
		handler, err := m.GetFunction(sendPacket.TargetNodeFunction)
		if err != nil {
			log.Warning(err.Error())
			return
		}

		// handler.
		fmt.Println(handler)
		// packet.TargetNodeFunction
		tickerTimeout.Stop()

		req := new(function.Request)
		res := new(function.Response)

		handlerResponse, err := handler(req, res)
		fmt.Println("---- SHOULD SEND TO SELF -- CALL FUNC LOCALLY ----")
		fmt.Println(handlerResponse)
		fmt.Println(err)

		return
	}
	m.transport.Send(sendPacket)

	// transpo.Send(packet)
	<-tickerTimeout.C

}
