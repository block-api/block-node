package block

import (
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

func (bn *BlockNode) ProcessPocketDiscovery(pocket transporter.Pocket[transporter.PayloadDiscovery]) {
	if pocket.Payload.Event == transporter.EventConnected {
		// Add to traffic Manager
		processDiscoveryBlocks(bn, pocket)
		return
	}

	if pocket.Payload.Event == transporter.EventDisconnected {
		bn.trafficManager.RemoveNodeID(pocket.FromID)
		return
	}

	if pocket.Payload.Event == transporter.EventHeartbeat {
		var knownNode bool = false
		knownNodes := bn.trafficManager.Nodes()

		for nodeID := range knownNodes {
			if pocket.FromID == nodeID {
				knownNode = true
				bn.trafficManager.UpdateLastSeen(nodeID)
				break
			}
		}

		if !knownNode {
			// Add to traffic Manager
			processDiscoveryBlocks(bn, pocket)
		}
		return
	}
}

func processDiscoveryBlocks(bn *BlockNode, pocket transporter.Pocket[transporter.PayloadDiscovery]) {
	for blockName, actions := range pocket.Payload.Blocks {
		err := bn.trafficManager.AddDestination(pocket.FromID, pocket.VersionName, blockName, actions)
		if err != nil {
			log.Warning(err.Error())
			continue
		}
	}
}
