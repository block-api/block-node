package traffic

import (
	"sync"

	"github.com/block-api/block-node/common/types"
)

type Destinations map[types.NodeVersionName]map[types.BlockName]map[types.ActionName][]types.NodeID
type Nodes []types.NodeID

// Manager struct
type Manager struct {
	destinationsMutex *sync.Mutex
	destinations      Destinations
	nodes             Nodes
}

// NewManager creates new Manager struct
func NewManager() Manager {
	return Manager{
		destinationsMutex: new(sync.Mutex),
		destinations:      make(map[types.NodeVersionName]map[types.BlockName]map[types.ActionName][]types.NodeID),
	}
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

	if nodeIDInSlice(m.nodes, nodeID) {
		return nil
	}

	if m.destinations[versionedName] == nil {
		m.destinations[versionedName] = make(map[types.BlockName]map[types.ActionName][]types.NodeID)
	}

	if m.destinations[versionedName][blockName] == nil {
		m.destinations[versionedName][blockName] = make(map[types.ActionName][]types.NodeID)
	}

	for _, actionName := range actions {
		m.destinations[versionedName][blockName][actionName] = append(m.destinations[versionedName][blockName][actionName], nodeID)
	}

	m.nodes = append(m.nodes, nodeID)

	return nil
}

// RemoveNodeID removes node id from all mentions in the topology
func (m *Manager) RemoveNodeID(nodeID types.NodeID) {
	m.destinationsMutex.Lock()
	defer m.destinationsMutex.Unlock()

	m.nodes = removeFromNodeIDSlice(m.nodes, nodeID)

	for nodeVersionName, blocks := range m.destinations {
		for blockName, actions := range blocks {
			for actionName, nodes := range actions {
				if nodeIDInSlice(nodes, nodeID) {
					m.destinations[nodeVersionName][blockName][actionName] = removeFromNodeIDSlice(nodes, nodeID)
				}
			}
		}
	}
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
