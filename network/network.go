package network

import (
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/transporter"
)

type Network struct {
	transporter transporter.Transporter
}

func (n *Network) Start() error {
	err := n.transporter.Subscribe(transporter.ChanDiscovery)
	if err != nil {
		log.Warning(err.Error())

		return err
	}

	return nil
}

func (n *Network) Send(pocket Pocket) {
	log.Debug("published to: " + string(pocket.Channel()))

	err := n.transporter.Send(pocket.Channel(), pocket)
	if err != nil {
		log.Warning(err.Error())
	}
}

func NewNetwork(transporter transporter.Transporter) Network {
	return Network{
		transporter: transporter,
	}
}
