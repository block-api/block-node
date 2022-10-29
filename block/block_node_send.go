package block

import (
	"encoding/json"
	"time"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

func (bn *BlockNode) Send(payload *transporter.PayloadMessage, targetAction *types.TargetAction) (*transporter.PayloadMessage, error) {
	var pocketPayload transporter.PayloadMessage

	if payload != nil {
		pocketPayload = *payload
	}

	pocket := transporter.NewPocket(transporter.ChanMessage, bn.nodeVersionName, bn.nodeID, nil, targetAction, pocketPayload)

	var err error

	if pocket.TargetAction != nil {
		err = pocket.TargetAction.Validate()
		if err != nil {
			return nil, err
		}

		if pocket.TargetAction.Name == types.NodeName(bn.options.Name) && pocket.TargetAction.Version == bn.options.Version && bn.blocks[pocket.TargetAction.Block] != nil {
			actions := bn.blocks[pocket.TargetAction.Block].Actions()

			if actions[pocket.TargetAction.Action] != nil {
				bn.sentHashes[pocket.Hash] = time.Now()

				payload, err := DecodePayload[transporter.PayloadMessage](pocket.Payload)
				if err != nil {
					log.Warning(err.Error())
					return nil, err
				}

				responsePayload, err := actions[pocket.TargetAction.Action](payload)
				if err != nil {
					return nil, err
				}

				return responsePayload, nil
			}
		}
		return nil, errors.ErrInvalidTargetAction
	}

	pocketBytes, err := json.Marshal(pocket)
	if err != nil {
		log.Warning(err.Error())
		return nil, err
	}

	err = bn.transporter.Send(pocket.Channel, pocketBytes)
	if err != nil {
		log.Warning(err.Error())
		return nil, err
	}

	bn.sentHashes[pocket.Hash] = time.Now()

	return nil, nil
}

func (bn *BlockNode) Response(payload *transporter.PayloadMessage) (*transporter.Pocket[[]byte], error) {
	var pocketPayload transporter.PayloadMessage

	if payload != nil {
		pocketPayload = *payload
	}

	pocket := transporter.NewPocket(transporter.ChanMessage, bn.nodeVersionName, bn.nodeID, nil, nil, pocketPayload)
	return &pocket, nil
}
