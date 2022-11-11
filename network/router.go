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
	"sync"

	"github.com/block-api/block-node/network/router"
	"github.com/block-api/block-node/params"
)

// knownNodesMaxSize max map in memory size to keep information about topology of known nodes
// it will try to find in memory, then in database
const knownNodesMaxSize = 50

type Router struct {
	config            *params.NetworkConfig
	knownNodes        map[string]router.Node
	knownNodesCounter int
	knownNodesLock    *sync.Mutex
}

func NewRouter(config *params.NetworkConfig) *Router {
	return &Router{
		config:            config,
		knownNodes:        make(map[string]router.Node),
		knownNodesCounter: 0,
		knownNodesLock:    new(sync.Mutex),
	}
}

func (r *Router) KnownNodes() *map[string]router.Node {
	return &r.knownNodes
}
