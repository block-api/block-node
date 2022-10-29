package block

import (
	"encoding/json"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

func (bn *BlockNode) Send(payload *transporter.PayloadMessage, targetAction *types.TargetAction) error {
	var pocketPayload transporter.PayloadMessage

	if payload != nil {
		pocketPayload = *payload
	}

	pocket := transporter.NewPocket(transporter.ChanMessage, bn.nodeVersionName, bn.nodeID, nil, targetAction, pocketPayload)

	// pocket transporter.Pocket[[]byte]
	var err error

	if pocket.TargetAction != nil {
		err = pocket.TargetAction.Validate()
		if err != nil {
			return err
		}

		if pocket.TargetAction.Name == types.NodeName(bn.options.Name) && pocket.TargetAction.Version == bn.options.Version && bn.blocks[pocket.TargetAction.Block] != nil {
			actions := bn.blocks[pocket.TargetAction.Block].Actions()

			if actions[pocket.TargetAction.Action] != nil {
				actions[pocket.TargetAction.Action](pocket.Payload)
				return nil
			}

		}

		return errors.ErrInvalidTargetAction
	}

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
