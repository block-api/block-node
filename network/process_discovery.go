package network

import (
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

func (n *Network) ProcessPocketDiscovery(pocket transporter.Pocket[transporter.PayloadDiscovery]) {
	if pocket.Payload.Event == transporter.EventConnected {
		// Add to traffic Manager
		for blockName, actions := range pocket.Payload.Blocks {
			err := n.trafficManager.AddDestination(pocket.FromID, pocket.VersionName, blockName, actions)
			if err != nil {
				log.Warning(err.Error())
				continue
			}
		}
		// return
	}

	if pocket.Payload.Event == transporter.EventDisconnected {
		n.trafficManager.RemoveNodeID(pocket.FromID)
		return
	}

	if pocket.Payload.Event == transporter.EventHeartbeat {
		var knownNode bool = false
		knownNodes := n.trafficManager.Nodes()

		for nodeID, _ := range knownNodes {
			if pocket.FromID == nodeID {
				knownNode = true
				n.trafficManager.UpdateLastSeen(nodeID)
				break
			}
		}

		if !knownNode {
			// Add to traffic Manager
			for blockName, actions := range pocket.Payload.Blocks {
				err := n.trafficManager.AddDestination(pocket.FromID, pocket.VersionName, blockName, actions)
				if err != nil {
					log.Warning(err.Error())
					continue
				}
			}
		}
		return
	}
}
