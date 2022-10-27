package block

import (
	"os"
	"strconv"
	"time"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/db"
	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network"
	"github.com/block-api/block-node/traffic"
	"github.com/block-api/block-node/transporter"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var instantiated bool

type BlockNode struct {
	nodeID            types.NodeID
	nodeVersionName   types.NodeVersionName
	blocks            map[types.BlockName]IBlock
	config            config.Config
	options           BlockNodeOptions
	transporter       transporter.Transporter
	network           network.Network
	database          db.Database
	trafficManager    traffic.Manager
	daemonChan        chan uint
	heartbeatInterval int
}

type BlockNodeOptions struct {
	Name    string
	Version uint
}

// Start will start BlockNode
func (bn *BlockNode) Start() {
	log.Debug("starting " + bn.options.Name + ", id: " + string(bn.nodeID))

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	configDir := os.Getenv("CONFIG_DIR")
	config := config.NewConfig()
	err = config.LoadFromFile(configDir + "/config.yml")

	if err != nil {
		panic(err)
	}

	bn.config = config

	err = bn.loadTransporter(config.Transporter)
	if err != nil {
		panic(err)
	}

	bn.loadDatabase()
	bn.loadNetwork()

	log.Default("# Name: " + string(bn.nodeVersionName) + " is running")
	log.Default("# NodeID: " + string(bn.nodeID))

	payload := transporter.PayloadDiscovery{
		Event:  transporter.EventConnected,
		Blocks: bn.Blocks(),
	}

	pocket := transporter.NewPocket[transporter.PayloadDiscovery](transporter.ChanDiscovery, bn.nodeVersionName, bn.nodeID, "", payload)
	bn.network.Send(pocket)

	go bn.daemon(bn.daemonChan)

}

// AddBlock adds new Block struct to BlockNode blocks map
func (bn *BlockNode) AddBlock(blocks ...IBlock) error {
	var bk IBlock

	for _, b := range blocks {
		if bn.blocks[b.GetName()] != bk {
			return errors.New(errors.ErrBlockAdded)
		}

		bn.blocks[b.GetName()] = b
		log.Debug("block added: " + b.GetName().String())

		bn.trafficManager.AddDestination(types.NodeID(bn.nodeID), types.NodeVersionName(bn.nodeVersionName), b.GetName(), b.ActionsNames())
	}

	return nil
}

func (bn BlockNode) Blocks() map[types.BlockName][]types.ActionName {
	var blocks map[types.BlockName][]types.ActionName = make(map[types.BlockName][]types.ActionName)

	for name, blck := range bn.blocks {
		actions := blck.Actions()

		for actionName := range actions {
			blocks[name] = append(blocks[name], actionName)
		}

	}
	return blocks
}

func (bn *BlockNode) Stop() error {
	payload := transporter.PayloadDiscovery{
		Event: transporter.EventDisconnected,
	}

	discoveryDisconnect := transporter.NewPocket(transporter.ChanDiscovery, bn.nodeVersionName, bn.nodeID, "", payload)
	bn.network.Send(discoveryDisconnect)

	bn.daemonChan <- 1

	err := bn.network.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (bn *BlockNode) NodeID() types.NodeID {
	return bn.nodeID
}

func (bn *BlockNode) GetName() string {
	return bn.options.Name
}

func (bn *BlockNode) Version() uint {
	return bn.options.Version
}

func (bn *BlockNode) VersionName() types.NodeVersionName {
	return bn.nodeVersionName
}

func (bn *BlockNode) Database() *db.Database {
	return &bn.database
}

func (bn *BlockNode) Network() *network.Network {
	return &bn.network
}

func (bn *BlockNode) loadDatabase() {
	bn.database = db.NewDatabase(&bn.config.Database)
}

func (bn *BlockNode) loadNetwork() {
	bn.network = network.NewNetwork(bn.nodeID, bn.transporter, &bn.trafficManager, &bn.database)
	err := bn.network.Start()

	if err != nil {
		log.Panic(err.Error())
	}
}

func (bn *BlockNode) daemon(daemonChan chan uint) {
	log.Debug("BlockNode daemon start")

	ticker := time.NewTicker(time.Duration(bn.heartbeatInterval) * time.Second)
L:
	for {
		select {
		case <-ticker.C:
			// send heartbeat
			payload := transporter.PayloadDiscovery{
				Event:  transporter.EventHeartbeat,
				Blocks: bn.Blocks(),
			}

			pocket := transporter.NewPocket(transporter.ChanDiscovery, bn.nodeVersionName, bn.nodeID, "", payload)
			bn.network.Send(pocket)

			for nodeID, lastSeen := range bn.trafficManager.Nodes() {
				if nodeID == bn.nodeID {
					continue
				}

				dateDiff := time.Since(lastSeen).Seconds()
				if dateDiff > float64(bn.heartbeatInterval+1) {
					bn.trafficManager.RemoveNodeID(nodeID)
				}
			}

		case <-daemonChan:
			break L
		}
	}

	log.Debug("BlockNode daemon quit")
}

var heartbeatInterval = 5

// NewBlockNode creates new BlockNode struct
func NewBlockNode(options *BlockNodeOptions) BlockNode {
	if instantiated {
		panic(errors.ErrBlockNodeInstantiated)
	}

	nodeID := types.NodeID("v" + strconv.Itoa(int(options.Version)) + "." + options.Name + "." + uuid.NewString())
	nodeVersionName := types.NodeVersionName("v" + strconv.Itoa(int(options.Version)) + "." + options.Name)

	bn := BlockNode{
		nodeID:            nodeID,
		nodeVersionName:   nodeVersionName,
		options:           *options,
		blocks:            make(map[types.BlockName]IBlock),
		transporter:       nil,
		trafficManager:    traffic.NewManager(nodeID),
		daemonChan:        make(chan uint),
		heartbeatInterval: heartbeatInterval,
	}

	instantiated = true

	return bn
}
