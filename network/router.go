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

	"github.com/block-api/block-node/block/sys/model"
	"github.com/block-api/block-node/block/sys/repo"
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
// const knownNodesMaxMemorySize = 50

type Router struct {
	knownNodesRepo    *repo.KnownNodeRepo
	nodeName          string
	nodeVersion       int
	config            *params.NetworkConfig
	knownNodes        map[string]*router.Node
	knownNodesCounter int
	knownNodesLock    *sync.Mutex
}

func NewRouter(nodeName string, nodeVersion int, config *params.NetworkConfig) *Router {
	rtr := &Router{
		knownNodesRepo:    repo.GetKnownNode(),
		nodeName:          nodeName,
		nodeVersion:       nodeVersion,
		config:            config,
		knownNodes:        make(map[string]*router.Node),
		knownNodesCounter: 0,
		knownNodesLock:    new(sync.Mutex),
	}

	dbNodesCount := rtr.knownNodesRepo.Count()
	if dbNodesCount > 0 {

	}

	if len(config.Nodes) > 0 {
		for _, knownNode := range config.Nodes {
			rtr.Add(knownNode.NodeID, &router.Node{
				Transport:  knownNode.Transport,
				NodeID:     knownNode.NodeID,
				PublicHost: knownNode.PublicHost,
				PublicPort: knownNode.PublicPort,
				Functions:  make(map[string]bool),
			})
		}

	}

	// fmt.Println(rtr.knownNodesRepo.Count())
	// counterRepo := repo.GetCounter()
	// resCounter, err := counterRepo.Get(params.DBSysKnownNodes)
	// if err != nil {
	// 	log.Warning(err.Error())
	// }
	// fmt.Println(resCounter)
	// fmt.Println(counterRepo)
	// rtr.knownNodesRepo.Get()
	return rtr
}

func (r *Router) Add(nodeID string, node *router.Node) error {
	r.knownNodesLock.Lock()
	defer r.knownNodesLock.Unlock()

	hasNode, err := r.knownNodesRepo.Has(nodeID)
	if err != nil {
		return err
	}

	if !hasNode {
		err := r.knownNodesRepo.Add(nodeID, model.KnownNode{
			Transport:  node.Transport,
			PublicHost: node.PublicHost,
			PublicPort: node.PublicPort,
			Functions:  node.Functions,
			LastSeen:   node.LastSeen,
		})
		if err != nil {
			return err
		}

		return nil
	}

	return ErrRouterNodeAlreadyExist
}

func (r *Router) GetNode(nodeID string) *router.Node {
	node, err := r.knownNodesRepo.Get(nodeID)
	if err != nil {
		return nil
	}

	return &router.Node{Transport: node.Transport, NodeID: nodeID, PublicHost: node.PublicHost, PublicPort: node.PublicPort, Functions: node.Functions, LastSeen: node.LastSeen}
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
