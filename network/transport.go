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
	"github.com/block-api/block-node/network/transport"
	"github.com/block-api/block-node/params"
)

type ITransport interface {
	Start(cReceive chan<- Packet) error
	Stop() error
	Send(packet Packet)
}

type Transport struct {
	transpo ITransport
}

func NewTransport[T TCPTransport | RedisTransport](networkManager *Manager, config *params.NetworkConfig) Transport {
	t := Transport{}
	switch config.Transport {
	case transport.TCP:
		t.transpo = TCPTransport{
			cReceive:       make(chan Packet),
			networkManager: networkManager,
		}

	case transport.REDIS:
		t.transpo = RedisTransport{
			cReceive:       make(chan Packet),
			networkManager: networkManager,
		}
	}
	return t
}

func (t *Transport) Start(cReceive chan<- Packet) error {
	return t.transpo.Start(cReceive)
}

func (t *Transport) Send(packet Packet) {
	t.transpo.Send(packet)
}
