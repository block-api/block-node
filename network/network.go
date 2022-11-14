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
	"errors"
	"fmt"
	"sync"

	"github.com/block-api/block-node/block/function"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network/transport"
	"github.com/block-api/block-node/params"
)

var (
	ErrAlreadyInstantiatied = errors.New("network manager is already instantiated")
	ErrConfigFileNotFound   = errors.New("config file not found")
)

var manager *Manager
var managerLock = new(sync.Mutex)

type Manager struct {
	nodeID               string
	config               *params.NetworkConfig
	functionManager      *function.Manager
	transport            Transport
	router               *Router
	cStopNetworkWorker   chan int
	cStopSenderkWorker   chan int
	cStopReceiverkWorker chan int
	cSend                chan Packet
	cReceive             chan Packet
	wgStop               *sync.WaitGroup
	cSendLock            *sync.Mutex
}

func NewManager(nodeID string, nodeName string, nodeVersion int, config *params.NetworkConfig, functionManager *function.Manager) (*Manager, error) {
	if manager == nil {
		managerLock.Lock()
		defer managerLock.Unlock()

		manager = &Manager{
			nodeID:               nodeID,
			config:               config,
			functionManager:      functionManager,
			router:               NewRouter(nodeName, nodeVersion, config),
			cStopNetworkWorker:   make(chan int),
			cStopSenderkWorker:   make(chan int),
			cStopReceiverkWorker: make(chan int),
			cSend:                make(chan Packet, 10),
			cReceive:             make(chan Packet),
			wgStop:               new(sync.WaitGroup),
			cSendLock:            new(sync.Mutex),
		}

		switch config.Transport {
		case transport.TCP:
			manager.transport = NewTransport[TCPTransport](manager, config)

		case transport.REDIS:
			manager.transport = NewTransport[RedisTransport](manager, config)
		}

		return manager, nil
	}

	return nil, ErrAlreadyInstantiatied
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) Router() *Router {
	return m.router
}

func (m *Manager) Start() error {
	// m.wgStop.Add(1)

	go m.networkWorker(m.cStopNetworkWorker, m.wgStop)

	// m.wgStop.Add(1)

	go m.senderWorker(m.nodeID, m.config, &m.transport, m.cStopSenderkWorker, m.cSend, m.wgStop)

	// m.wgStop.Add(1)

	go m.receiverWorker(m.transport, m.cStopReceiverkWorker, m.cReceive, m.wgStop)

	return nil
}

func (m *Manager) Stop() error {

	m.cStopSenderkWorker <- 1
	m.cStopReceiverkWorker <- 1
	m.cStopNetworkWorker <- 1

	// m.wgStop.Wait()

	return nil
}

func (m *Manager) GetFunction(name string) (function.Handler, error) {
	return m.functionManager.Get(name)
}

func (n *Manager) Send(packet Packet) error {
	err := packet.Validate()
	if err != nil {

		return err
	}

	log.Default("### SEND ###")
	fmt.Println(packet)
	n.cSendLock.Lock()
	n.cSend <- packet
	n.cSendLock.Unlock()
	return nil
}
