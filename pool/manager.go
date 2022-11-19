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

// Package pool
package pool

import (
	"errors"
	"sync"
)

var (
	manager     *Manager
	managerLock = new(sync.Mutex)
)
var (
	ErrManagerInstantiated = errors.New("manager is already instantiated")
	ErrPoolExist           = errors.New("pool already exist")
	ErrPoolNotExist        = errors.New("pool doest not exist")
)

type Manager struct {
	pools map[string]Pool
	lock  *sync.Mutex
}

func NewManager() (*Manager, error) {
	if manager == nil {
		managerLock.Lock()
		defer managerLock.Unlock()

		manager = &Manager{
			pools: make(map[string]Pool),
			lock:  new(sync.Mutex),
		}

		return manager, nil
	}

	return nil, ErrManagerInstantiated
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) NewPool(name string, maxWorkers int, workerHandler WorkerHandler) error {
	if m.pools[name].name != "" {
		return ErrPoolExist
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.pools[name] = Pool{
		maxWorkers: maxWorkers,
		worker:     workerHandler,
		cIn:        make(chan interface{}),
		cOut:       make(chan interface{}),
		lockIn:     new(sync.Mutex),
	}

	return nil
}

func (m *Manager) SendToPool(poolName string, payload interface{}) error {
	if m.pools[poolName].name == "" {
		return ErrPoolNotExist
	}

	m.pools[poolName].lockIn.Lock()
	defer m.pools[poolName].lockIn.Unlock()

	m.pools[poolName].cIn <- payload

	return nil
}

func (m *Manager) ListenForResults(poolName string, callback func(payload interface{})) error {
	if m.pools[poolName].name == "" {
		return ErrPoolNotExist
	}

	go func(cOut <-chan interface{}, callback func(payload interface{})) {
		for {
			workerRes := <-cOut
			callback(workerRes)
		}
	}(m.pools[poolName].cOut, callback)

	return nil
}
