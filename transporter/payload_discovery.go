package transporter

import "github.com/block-api/block-node/common/types"

type PayloadDiscovery struct {
	NodeID  string                                 `json:"node_id"`
	Name    string                                 `json:"name"`
	Version uint                                   `json:"version"`
	Blocks  map[types.BlockName][]types.ActionName `json:"blocks"`
	Event   Event                                  `json:"event"`
}
