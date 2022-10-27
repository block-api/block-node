package network

import (
	"fmt"

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
	}

	fmt.Println(n.trafficManager.Destinations())
}
