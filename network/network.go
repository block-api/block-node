package network

import (
	"encoding/json"

	"github.com/block-api/block-node/db"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

type Network struct {
	db          *db.Database
	transporter transporter.Transporter
}

func (n *Network) Start() error {
	err := n.transporter.Subscribe(transporter.ChanDiscovery, n.Receive)

	if err != nil {
		log.Warning(err.Error())

		return err
	}

	return nil
}

func (n *Network) Send(pocket transporter.Pocket[[]byte]) {
	pocketBytes, err := json.Marshal(pocket)
	if err != nil {
		log.Warning(err.Error())
	}

	err = n.transporter.Send(pocket.Channel, pocketBytes)
	if err != nil {
		log.Warning(err.Error())
	}
}

func (n *Network) Receive(pocket transporter.Pocket[[]byte]) {
	log.Debug("Network.Receive Pocket[[]byte]")

	if pocket.Channel == transporter.ChanDiscovery {
		discoveryPayload, err := decodeDiscoveryPayload(&pocket)
		if err != nil {
			log.Warning(err.Error())
		}

		newPocket := transporter.Pocket[transporter.PayloadDiscovery]{
			Channel: pocket.Channel,
			From:    pocket.From,
			Target:  pocket.Target,
			Payload: discoveryPayload,
			Hash:    pocket.Hash,
		}

		n.ProcessPocketDiscovery(newPocket)
	}
}

func NewNetwork(transporter transporter.Transporter, db *db.Database) Network {
	return Network{
		db:          db,
		transporter: transporter,
	}
}
