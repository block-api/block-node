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
