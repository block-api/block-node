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
	"math/rand"
	"sync"
	"time"

	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network/delivery"
	"github.com/block-api/block-node/network/router"
	"github.com/block-api/block-node/params"
)

var (
	ErrRouterUnknownNode      = errors.New("unknown router node id")
	ErrRouterNodeAlreadyExist = errors.New("router node id already exist")
)

// knownNodesMaxSize max map in memory size to keep information about topology of known nodes
// it will try to find in memory, then in database
const knownNodesMaxSize = 50

type Router struct {
	nodeName          string
	nodeVersion       int
	config            *params.NetworkConfig
	knownNodes        map[string]*router.Node
	knownNodesCounter int
	knownNodesLock    *sync.Mutex
}

func NewRouter(nodeName string, nodeVersion int, config *params.NetworkConfig) *Router {
	return &Router{
		nodeName:          nodeName,
		nodeVersion:       nodeVersion,
		config:            config,
		knownNodes:        make(map[string]*router.Node),
		knownNodesCounter: 0,
		knownNodesLock:    new(sync.Mutex),
	}
}

func (r *Router) Add(nodeID string, node *router.Node) error {
	r.knownNodesLock.Lock()
	defer r.knownNodesLock.Unlock()

	if r.knownNodes[nodeID] != nil && r.knownNodes[nodeID].LastSeen >= node.LastSeen {
		log.Warning(ErrRouterNodeAlreadyExist.Error())
		return ErrRouterNodeAlreadyExist
	}

	r.knownNodes[nodeID] = node

	return nil
}

func (r *Router) GetNode(nodeID string) *router.Node {
	return r.knownNodes[nodeID]
}

func (r *Router) Remove(nodeID string) error {
	if r.knownNodes[nodeID] == nil {
		return ErrRouterUnknownNode
	}

	r.knownNodesLock.Lock()
	defer r.knownNodesLock.Unlock()

	delete(r.knownNodes, nodeID)

	return nil
}

func (r *Router) UpdateLastSeen(knownNodeID string, lastSeenAt int64) {
	r.knownNodesLock.Lock()
	defer r.knownNodesLock.Unlock()

	r.knownNodes[knownNodeID].LastSeen = lastSeenAt
}

func (r *Router) KnownNodes() *map[string]*router.Node {
	return &r.knownNodes
}

func (r *Router) GetTarget(deliveryMethod delivery.Type, targetFunction string, targetNodeID string) []*router.Node {
	var targetNodes []*router.Node = make([]*router.Node, 0)

	if deliveryMethod == delivery.All {
		for _, node := range r.knownNodes {
			targetNodes = append(targetNodes, node)
		}

		return targetNodes

	}

	if deliveryMethod == delivery.Random {

		if len(r.knownNodes) > 0 {
			var node *router.Node
			var nodes = make([]*router.Node, 0)

			for _, node := range r.knownNodes {
				nodes = append(nodes, node)
			}

			rand.Seed(time.Now().Unix())

			node = nodes[rand.Intn(len(nodes))]
			targetNodes = append(targetNodes, node)

			return targetNodes
		}
	}
	return targetNodes
}
