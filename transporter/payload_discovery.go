package transporter

import "github.com/block-api/block-node/utils"

type PayloadDiscovery struct {
	NodeID  string                                 `json:"node_id"`
	Name    string                                 `json:"name"`
	Version uint                                   `json:"version"`
	Blocks  map[utils.BlockName][]utils.ActionName `json:"blocks"`
	Event   Event                                  `json:"event"`
}
