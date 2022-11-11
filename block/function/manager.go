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
package function

import (
	"errors"
)

var (
	ErrFunctionAlreadyExist = errors.New("function already exist")
	ErrFunctionDoesNotExist = errors.New("function does not exist")
)

type Manager struct {
	functions map[string]Handler
}

func NewManager() *Manager {
	manager := &Manager{
		functions: make(map[string]Handler),
	}

	manager.Add("sys.status", SysStatusFunction)

	return manager
}

func (m *Manager) Add(name string, fn Handler) error {
	if m.functions[name] != nil {
		return ErrFunctionAlreadyExist
	}

	m.functions[name] = fn

	return nil
}

func (m *Manager) Get(name string) (Handler, error) {
	if m.functions[name] == nil {
		return nil, ErrFunctionDoesNotExist
	}

	return m.functions[name], nil
}
