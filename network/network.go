package network

import (
	"encoding/json"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/db"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/traffic"
	"github.com/block-api/block-node/transporter"
)

type Network struct {
	nodeID         types.NodeID
	db             *db.Database
	transporter    transporter.Transporter
	trafficManager *traffic.Manager
}

func (n *Network) Start() error {
	err := n.transporter.Subscribe(transporter.ChanDiscovery, n.Receive)

	if err != nil {
		log.Warning(err.Error())

		return err
	}

	return nil
}

func (n *Network) Stop() error {
	n.trafficManager.Stop()
	err := n.transporter.Disconnect()

	if err != nil {
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

func (n *Network) Receive(payload []byte) {
	log.Debug("Network.Receive []byte")

	pocket, err := decodePocket(payload)
	if err != nil {
		log.Warning(err.Error())
		return
	}

	if (pocket.FromID == n.nodeID) || (pocket.TargetID != "" && pocket.TargetID != n.nodeID) {
		log.Debug("skip self/not target")
		return
	}

	if pocket.Channel == transporter.ChanDiscovery {
		discoveryPayload, err := decodePayload[transporter.PayloadDiscovery](&pocket)
		if err != nil {
			log.Warning(err.Error())
			return
		}

		newPocket := transporter.Pocket[transporter.PayloadDiscovery]{
			Channel:     pocket.Channel,
			VersionName: pocket.VersionName,
			FromID:      pocket.FromID,
			TargetID:    pocket.TargetID,
			Payload:     discoveryPayload,
			Hash:        pocket.Hash,
		}

		n.ProcessPocketDiscovery(newPocket)
	}
}

func NewNetwork(nodeID types.NodeID, transporter transporter.Transporter, trafficManager *traffic.Manager, db *db.Database) Network {
	return Network{
		nodeID:         nodeID,
		db:             db,
		transporter:    transporter,
		trafficManager: trafficManager,
	}
}
