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
package block

import (
	"errors"
	"sync"

	"github.com/block-api/block-node/params"
)

var (
	ErrAlreadyInstantiatied = errors.New("node is already instantiated")
)

var (
	node     *Node
	nodeLock = new(sync.Mutex)
)

// Node is main struct describing node
type Node struct {
	// id is unique identifier for node instance
	id     string
	config *params.NodeConfig
}

// New creates new node instance, there can be only one instance of node in your program
func NewNode(config *params.NodeConfig) (*Node, error) {
	if node == nil {
		nodeLock.Lock()
		defer nodeLock.Unlock()

		node = &Node{
			config: config,
		}
		return node, nil
	}
	return nil, ErrAlreadyInstantiatied
}

// Get pointer to node
func GetNode() *Node {
	return node
}

// ID returns node unique identifier
func (n *Node) ID() string {
	return n.id
}
