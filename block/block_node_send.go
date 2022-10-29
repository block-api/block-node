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
				payload, err := DecodePayload[transporter.PayloadMessage](pocket.Payload)
				if err != nil {
					log.Warning(err.Error())
					return nil, err
				}

				// local invocation
				responsePayload, err := actions[pocket.TargetAction.Action](payload)
				if err != nil {
					return nil, err
				}

				return responsePayload, nil
			}

			return nil, errors.ErrInvalidTargetAction
		}

		if bn.trafficManager.DestinationExist(*pocket.TargetAction) {
			pocketBytes, err := json.Marshal(pocket)
			if err != nil {
				log.Warning(err.Error())
				return nil, err
			}

			// TODO: get random destination id from TrafficManager
			// in the future other options eg. lowest CPU usage

			err = bn.transporter.Send(pocket.Channel, pocketBytes)
			if err != nil {
				log.Warning(err.Error())
				return nil, err
			}

			bn.sentHashesMutex.Lock()
			bn.sentHashes[pocket.Hash] = &SentHash{time: time.Now(), responseChan: make(chan transporter.Pocket[[]byte])}

			// wait for X (timeout) seconds for response
			var responsePayloadChan chan *transporter.PayloadMessage = make(chan *transporter.PayloadMessage)

			go func(responseChan chan *transporter.PayloadMessage, receivedResponsePocket chan transporter.Pocket[[]byte]) {
				timer := time.NewTimer(15 * time.Second)
			L:
				for {
					select {
					case res := <-receivedResponsePocket:
						receivedPayload, err := DecodePayload[transporter.PayloadMessage](res.Payload)
						if err != nil {
							log.Warning(err.Error())
						}

						responseChan <- &receivedPayload
					case <-timer.C:
						responseChan <- &transporter.PayloadMessage{}
						break L
					}
				}
			}(responsePayloadChan, bn.sentHashes[pocket.Hash].responseChan)
			bn.sentHashesMutex.Unlock()

			receivedPayload := <-responsePayloadChan

			bn.sentHashesMutex.Lock()
			delete(bn.sentHashes, pocket.Hash)
			bn.sentHashesMutex.Unlock()

			return receivedPayload, nil
		}

		return nil, errors.ErrInvalidTargetAction
	}

	return nil, errors.ErrInvalidTargetAction
}
