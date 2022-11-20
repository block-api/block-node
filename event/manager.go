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

// Package event
package event

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyInstantiated = errors.New("event manager already instantiated")
	ErrNameDoesNotExist    = errors.New("event name does not exist")
)

var (
	manager     *Manager
	managerLock = new(sync.Mutex)
)

type Manager struct {
	eventsCallbacks map[string][]Listener
	lock            *sync.RWMutex
}

func NewManager() (*Manager, error) {
	if manager == nil {
		managerLock.Lock()
		defer managerLock.Unlock()

		manager = &Manager{
			eventsCallbacks: make(map[string][]Listener),
			lock:            new(sync.RWMutex),
		}
		return manager, nil
	}
	return nil, ErrAlreadyInstantiated
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) On(eventName string, callback Listener) int {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.eventsCallbacks[eventName] == nil {
		m.eventsCallbacks[eventName] = make([]Listener, 0)
	}

	m.eventsCallbacks[eventName] = append(m.eventsCallbacks[eventName], callback)

	return len(m.eventsCallbacks[eventName]) - 1
}

func (m *Manager) Emit(e Event) error {
	if m.eventsCallbacks[string(e.Name)] == nil {
		return ErrNameDoesNotExist
	}

	go func(e Event, listeners []Listener) {
		for _, listener := range listeners {
			listener(e)
		}
	}(e, m.eventsCallbacks[string(e.Name)])
	return nil
}
