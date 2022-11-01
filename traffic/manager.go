package traffic

import (
	"github.com/block-api/block-node/common"
	"github.com/block-api/block-node/transporter"
	"math/rand"
	"sync"
	"time"

	"github.com/block-api/block-node/common/types"
)

type Destinations map[types.NodeVersionName]map[types.BlockName]map[types.ActionName][]types.NodeID
type Nodes map[types.NodeID]time.Time

// Manager struct
type Manager struct {
	nodeID            types.NodeID
	destinationsMutex *sync.Mutex
	destinations      Destinations
	nodes             Nodes
}

// NewManager creates new Manager struct
func NewManager(nodeID *types.NodeID) Manager {
	manager := Manager{
		nodeID:            *nodeID,
		destinationsMutex: new(sync.Mutex),
		destinations:      make(map[types.NodeVersionName]map[types.BlockName]map[types.ActionName][]types.NodeID),
		nodes:             make(Nodes),
	}

	return manager
}

// Nodes returns array of known node ID's
func (m *Manager) Nodes() Nodes {
	return m.nodes
}

// Destinations returns network topology
func (m *Manager) Destinations() Destinations {
	return m.destinations
}

// DestinationExist return true if target action is found in the network
func (m *Manager) DestinationExist(targetAction types.TargetAction) bool {
	nodeVersionName := common.CreateNodeVersionName(targetAction.Version, string(targetAction.Name))

	if m.destinations[nodeVersionName] == nil && m.destinations[nodeVersionName][targetAction.Block] == nil && m.destinations[nodeVersionName][targetAction.Block][targetAction.Action] == nil {
		return false
	}

	return true
}

// GetDeliveryTargetNodeID return NodeID for target action based on set deliveryMethod in config file
func (m *Manager) GetDeliveryTargetNodeID(deliveryMethod types.DeliveryMethod, targetAction types.TargetAction) (*types.NodeID, error) {
	nodeVersionName := common.CreateNodeVersionName(targetAction.Version, string(targetAction.Name))

	if deliveryMethod == transporter.RandomDelivery {
		availableNodes := m.destinations[nodeVersionName][targetAction.Block][targetAction.Action]
		if len(availableNodes) < 1 {
			return nil, ErrNoNodeFound
		}

		rand.Seed(time.Now().Unix())

		return &availableNodes[rand.Intn(len(availableNodes))], nil
	}

	return nil, transporter.ErrInvalidDeliveryMethod
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

func (m *Manager) UpdateLastSeen(nodeID types.NodeID) {
	m.nodes[nodeID] = time.Now()
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
