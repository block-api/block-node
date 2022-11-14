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
	"github.com/block-api/block-node/log"
)

type RedisTransport struct {
	ITransport
	networkManager *Manager
	cReceive       chan Packet
}

func (rt RedisTransport) Start(cReceive chan<- Packet) error {
	log.Debug("-- redis_transport::start --")
	return nil
}

func (rt RedisTransport) Stop() error {
	log.Debug("-- redis_transport::stop --")
	return nil
}

func (rt RedisTransport) Send(packet Packet) {
	log.Debug("-- redis_transport::send --")
}
