package traffic

import (
	"sync"

	"github.com/block-api/block-node/utils"
)

type Destinations map[utils.NodeVersionName]map[utils.BlockName]map[utils.ActionName][]utils.NodeID
type Nodes []utils.NodeID

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
		destinations:      make(map[utils.NodeVersionName]map[utils.BlockName]map[utils.ActionName][]utils.NodeID),
	}
}

// Nodes returns array of known node ID's
func (m *Manager) Nodes() Nodes {
	return m.nodes
}

// AddDestination adds information about topology of known nodes to manager
func (m *Manager) AddDestination(nodeID utils.NodeID, versionedName utils.NodeVersionName, blockName utils.BlockName, actions []utils.ActionName) error {
	m.destinationsMutex.Lock()
	defer m.destinationsMutex.Unlock()

	if m.destinations[versionedName] == nil {
		m.destinations[versionedName] = make(map[utils.BlockName]map[utils.ActionName][]utils.NodeID)
	}

	if m.destinations[versionedName][blockName] == nil {
		m.destinations[versionedName][blockName] = make(map[utils.ActionName][]utils.NodeID)
	}

	for _, actionName := range actions {
		m.destinations[versionedName][blockName][actionName] = append(m.destinations[versionedName][blockName][actionName], nodeID)
	}

	return nil
}

func (m *Manager) Send(data any) error {

	return nil
}
