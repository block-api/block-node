package transporter

import "github.com/block-api/block-node/common/types"

type PayloadDiscovery struct {
	Blocks map[types.BlockName][]types.ActionName `json:"blocks"`
	Event  Event                                  `json:"event"`
}
