package block

import (
	"os"
	"strconv"

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
	nodeID          types.NodeID
	nodeVersionName types.NodeVersionName
	blocks          map[types.BlockName]IBlock
	config          config.Config
	options         BlockNodeOptions
	transporter     transporter.Transporter
	network         network.Network
	database        db.Database
	trafficManager  traffic.Manager
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
}

// AddBlock adds new Block struct to BlockNode blocks map
func (bn *BlockNode) AddBlock(blocks ...IBlock) error {
	var bk IBlock

	for _, b := range blocks {
		// fmt.Println(b.GetName())
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
	err := bn.transporter.Disconnect()

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
	bn.network = network.NewNetwork(bn.nodeID, bn.transporter, &bn.database)
	err := bn.network.Start()

	if err != nil {
		log.Panic(err.Error())
	}
}

// NewBlockNode creates new BlockNode struct
func NewBlockNode(options *BlockNodeOptions) BlockNode {
	if instantiated {
		panic(errors.ErrBlockNodeInstantiated)
	}

	nodeID := "v" + strconv.Itoa(int(options.Version)) + "." + options.Name + "." + uuid.NewString()
	nodeVersionName := "v" + strconv.Itoa(int(options.Version)) + "." + options.Name

	bn := BlockNode{
		nodeID:          types.NodeID(nodeID),
		nodeVersionName: types.NodeVersionName(nodeVersionName),
		options:         *options,
		blocks:          make(map[types.BlockName]IBlock),
		transporter:     nil,
		trafficManager:  traffic.NewManager(),
	}

	instantiated = true

	return bn
}
