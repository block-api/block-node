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
	"os"
	"sync"

	"github.com/block-api/block-node/common"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network"
	"github.com/block-api/block-node/params"
	"github.com/joho/godotenv"
)

var (
	ErrAlreadyInstantiatied = errors.New("node is already instantiated")
	ErrConfigFileNotFound   = errors.New("config file not found")
)

var (
	node     *Node
	nodeLock = new(sync.Mutex)
)

// Node is main struct describing node
type Node struct {
	// id is unique identifier for node instance
	id             string
	config         *params.NodeConfig
	networkManager *network.Manager
	cStop          chan int
	wgNodeWorker   *sync.WaitGroup
}

// New creates new node instance, there can be only one instance of node in your program
func NewNode() (*Node, error) {
	if node == nil {
		nodeLock.Lock()
		defer nodeLock.Unlock()

		_ = godotenv.Load()

		configFile := os.Getenv("CONFIG_FILE")
		dataDir := os.Getenv("DATA_DIR")

		if dataDir == "" {
			dataDir = params.DefaultDataDir
		}

		if configFile == "" {
			configFile = dataDir + "/config.yml"
		}

		var config = new(params.NodeConfig)

		cfgFile, err := common.OpenFile(configFile, common.YML)
		if err != nil {
			log.Fatal(err.Error())
		}

		if cfgFile != nil {
			err = cfgFile.Parse(config)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		// todo: config validation

		networkManager, err := network.NewManager(&config.Network)
		if err != nil {
			log.Fatal(err.Error())
		}

		node = &Node{
			config:         config,
			networkManager: networkManager,
			cStop:          make(chan int),
			wgNodeWorker:   new(sync.WaitGroup),
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

// Config returns pointer to NodeConfig
func (n *Node) Config() *params.NodeConfig {
	return n.config
}

// Stop sends information to cStop channel to stop program
func (n *Node) Stop() {
	n.cStop <- 1

	n.wgNodeWorker.Wait()

	log.Debug("node stopped")
}

// Start runs all needed actions to start node
func (n *Node) Start() {
	n.wgNodeWorker.Add(1)

	go nodeWorker(n)
}
