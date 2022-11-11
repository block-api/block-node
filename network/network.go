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
	"sync"

	"github.com/block-api/block-node/block/function"
	"github.com/block-api/block-node/params"
)

var (
	ErrAlreadyInstantiatied = errors.New("network manager is already instantiated")
	ErrConfigFileNotFound   = errors.New("config file not found")
)

var manager *Manager
var managerLock = new(sync.Mutex)

type Manager struct {
	config          *params.NetworkConfig
	functionManager *function.Manager
}

func NewManager(config *params.NetworkConfig, functionManager *function.Manager) (*Manager, error) {
	if manager == nil {
		managerLock.Lock()
		defer managerLock.Unlock()

		manager = &Manager{
			config:          config,
			functionManager: functionManager,
		}

		return manager, nil
	}

	return nil, ErrAlreadyInstantiatied
}

func (m *Manager) GetFunction(name string) (function.Handler, error) {
	return m.functionManager.Get("sys.status")
}
