package network

import (
	"fmt"

	"github.com/block-api/block-node/transporter"
)

func (n *Network) ProcessPocketDiscovery(pocket transporter.Pocket[transporter.PayloadDiscovery]) {

	fmt.Println(pocket)
	// fmt.Println(pocket.Payload.Blocks)
	// db := n.db.GetLevelDB("discovery")
	// fmt.Println(db)
}
