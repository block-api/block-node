package block

import (
	"os"
	"strconv"

	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/db"
	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network"
	"github.com/block-api/block-node/transporter"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var instantiated bool

type BlockNode struct {
	nodeID      string
	blocks      map[string]IBlock
	config      config.Config
	options     BlockNodeOptions
	transporter transporter.Transporter
	network     network.Network
	database    db.Database
}

type BlockNodeOptions struct {
	Name    string
	Version uint
}

// Start will start BlockNode
func (bn *BlockNode) Start() {
	log.Debug("starting " + bn.options.Name)

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
}

// AddBlock adds new Block struct to BlockNode blocks map
func (bn *BlockNode) AddBlock(blocks ...IBlock) error {
	var bk IBlock

	for _, b := range blocks {
		if bn.blocks[b.GetName()] != bk {
			return errors.New(errors.ErrBlockAdded)
		}

		bn.blocks[b.GetName()] = b

		log.Debug("block added: " + b.GetName())
	}

	return nil
}

func (bn *BlockNode) Stop() error {
	err := bn.transporter.Disconnect()

	if err != nil {
		return err
	}

	return nil
}

func (bn *BlockNode) NodeID() string {
	return bn.nodeID
}

func (bn *BlockNode) loadDatabase() {
	bn.database = db.NewDatabase(&bn.config.Database)
}

func (bn *BlockNode) loadNetwork() {
	bn.network = network.NewNetwork(bn.transporter)
	err := bn.network.Start()

	if err != nil {
		log.Panic(err.Error())
	}

	bn.network.Send(transporter.ChanDiscovery, map[string]string{
		"status": "ok",
	})
}

// NewBlockNode creates new BlockNode struct
func NewBlockNode(options *BlockNodeOptions) BlockNode {
	if instantiated {
		panic(errors.ErrBlockNodeInstantiated)
	}

	nodeID := "v" + strconv.Itoa(int(options.Version)) + "." + options.Name + "." + uuid.NewString()

	bn := BlockNode{
		nodeID:      nodeID,
		options:     *options,
		blocks:      make(map[string]IBlock),
		transporter: nil,
	}

	instantiated = true

	return bn
}
