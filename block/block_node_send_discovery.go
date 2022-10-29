package block

import (
	"encoding/json"

	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

func (bn *BlockNode) sendDiscovery(payload *transporter.PayloadDiscovery) error {
	var pocketPayload transporter.PayloadDiscovery

	if payload != nil {
		pocketPayload = *payload
	}

	pocket := transporter.NewPocket(transporter.ChanDiscovery, bn.nodeVersionName, bn.nodeID, nil, nil, pocketPayload)

	pocketBytes, err := json.Marshal(pocket)
	if err != nil {
		log.Warning(err.Error())
		return err
	}

	err = bn.transporter.Send(pocket.Channel, pocketBytes)
	if err != nil {
		log.Warning(err.Error())
		return err
	}

	return nil
}
