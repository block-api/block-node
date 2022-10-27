package traffic

import (
	"fmt"
	"sync"
	"time"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/log"
)

var heartbeatInterval = 5

type Destinations map[types.NodeVersionName]map[types.BlockName]map[types.ActionName][]types.NodeID
type Nodes map[types.NodeID]time.Time

// Manager struct
type Manager struct {
	nodeID            types.NodeID
	destinationsMutex *sync.Mutex
	destinations      Destinations
	nodes             Nodes
	daemonChan        chan uint
}

// NewManager creates new Manager struct
func NewManager(nodeID types.NodeID) Manager {
	manager := Manager{
		nodeID:            nodeID,
		destinationsMutex: new(sync.Mutex),
		destinations:      make(map[types.NodeVersionName]map[types.BlockName]map[types.ActionName][]types.NodeID),
		nodes:             make(Nodes),
		daemonChan:        make(chan uint),
	}

	go manager.daemon(manager.daemonChan)

	return manager
}

func (m *Manager) Stop() {
	m.daemonChan <- 1
}

// Nodes returns array of known node ID's
func (m *Manager) Nodes() Nodes {
	return m.nodes
}

func (m *Manager) Destinations() Destinations {
	return m.destinations
}

// AddDestination adds information about topology of known nodes to manager
func (m *Manager) AddDestination(nodeID types.NodeID, versionedName types.NodeVersionName, blockName types.BlockName, actions []types.ActionName) error {
	m.destinationsMutex.Lock()
	defer m.destinationsMutex.Unlock()

	if m.destinations[versionedName] == nil {
		m.destinations[versionedName] = make(map[types.BlockName]map[types.ActionName][]types.NodeID)
	}

	if m.destinations[versionedName][blockName] == nil {
		m.destinations[versionedName][blockName] = make(map[types.ActionName][]types.NodeID)
	}

	for _, actionName := range actions {
		m.destinations[versionedName][blockName][actionName] = append(m.destinations[versionedName][blockName][actionName], nodeID)
	}

	m.nodes[nodeID] = time.Now()

	return nil
}

// RemoveNodeID removes node id from all mentions in the topology
func (m *Manager) RemoveNodeID(nodeID types.NodeID) {
	m.destinationsMutex.Lock()
	defer m.destinationsMutex.Unlock()

	for nodeVersionName, blocks := range m.destinations {
		for blockName, actions := range blocks {
			for actionName, nodes := range actions {
				if nodeIDInSlice(nodes, nodeID) {
					m.destinations[nodeVersionName][blockName][actionName] = removeFromNodeIDSlice(nodes, nodeID)

					if len(m.destinations[nodeVersionName][blockName][actionName]) < 1 {
						delete(m.destinations[nodeVersionName][blockName], actionName)
					}
				}
			}

			if len(m.destinations[nodeVersionName][blockName]) < 1 {
				delete(m.destinations[nodeVersionName], blockName)
			}
		}

		if len(m.destinations[nodeVersionName]) < 1 {
			delete(m.destinations, nodeVersionName)
		}
	}

	delete(m.nodes, nodeID)
}

func (m *Manager) daemon(daemonChan chan uint) {
	log.Debug("Traffic Manager daemon start")

	ticker := time.NewTicker(time.Duration(heartbeatInterval) * time.Second)
L:
	for {
		select {
		case <-ticker.C:
			for nodeID, lastSeen := range m.nodes {
				if nodeID == m.nodeID {
					continue
				}

				dateDiff := time.Since(lastSeen).Seconds()
				if dateDiff > float64(heartbeatInterval+1) {
					m.RemoveNodeID(nodeID)
				}
			}

			fmt.Println(m.destinations)
		case <-daemonChan:
			break L
		}
	}

	log.Debug("Traffic Manager daemon quit")
}

func nodeIDInSlice(slice []types.NodeID, nodeID types.NodeID) bool {
	for _, id := range slice {
		if id == nodeID {
			return true
		}
	}

	return false
}

func removeFromNodeIDSlice(slice []types.NodeID, value types.NodeID) []types.NodeID {
	var removeIDX int
	for idx, id := range slice {
		if id == value {
			removeIDX = idx
			break
		}
	}

	return append(slice[:removeIDX], slice[removeIDX+1:]...)
}
