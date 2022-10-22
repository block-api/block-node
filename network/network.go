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

func (n *Network) Send(channel transporter.Channel, payload interface{}) {
	log.Debug("published to: " + string(channel))

	err := n.transporter.Send(channel, payload)
	if err != nil {
		log.Warning(err.Error())
	}
}

func NewNetwork(transporter transporter.Transporter) Network {
	return Network{
		transporter: transporter,
	}
}
