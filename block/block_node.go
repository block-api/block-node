package block

import (
	"encoding/json"
	"github.com/block-api/block-node/common"
	"os"
	"sync"
	"time"

	"github.com/block-api/block-node/common/types"
	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/db"
	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/traffic"
	"github.com/block-api/block-node/transporter"
	"github.com/joho/godotenv"
)

var lock = new(sync.Mutex)
var blockNode *BlockNode

type SentHash struct {
	hash         string
	time         time.Time
	responseChan chan transporter.Pocket[[]byte]
}

type BlockNode struct {
	nodeID          types.NodeID
	nodeVersionName types.NodeVersionName
	blocks          map[types.BlockName]IBlock
	config          config.Config
	options         BlockNodeOptions
	transporter     transporter.Transporter
	database        *db.Database
	trafficManager  traffic.Manager
	daemonChan      chan uint
	sentHashes      map[string]*SentHash
	sentHashesMutex *sync.Mutex
}

type BlockNodeOptions struct {
	Name    string
	Version uint
}

func GetBlockNode() *BlockNode {
	return blockNode
}

// Start will start BlockNode
func (bn *BlockNode) Start() error {
	log.Debug("starting " + bn.options.Name + ", id: " + string(bn.nodeID))

	_ = godotenv.Load()

	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		panic(errors.ErrConfigDirMissing)
	}

	configFile := config.NewConfig()
	err := configFile.LoadFromFile(configDir + "/config.yml")
	if err != nil {
		panic(err)
	}

	if os.Getenv("DEBUG") == "" {
		if configFile.Debug {
			_ = os.Setenv("DEBUG", "true")
		} else {
			_ = os.Setenv("DEBUG", "false")
		}
	}

	bn.config = configFile
	err = bn.loadTransporter(bn.config.Transporter)
	if err != nil {
		panic(err)
	}

	if bn.config.Transporter.Timeout < 1 {
		bn.config.Transporter.Timeout = 30
	}

	if bn.config.Transporter.HeartbeatInterval < 1 {
		bn.config.Transporter.HeartbeatInterval = 5
	}

	bn.loadDatabase()

	err = bn.transporter.Subscribe(transporter.ChanDiscovery, bn.Receive)
	if err != nil {
		log.Warning(err.Error())
		return err
	}

	err = bn.transporter.Subscribe(transporter.ChanMessage, bn.Receive)
	if err != nil {
		log.Warning(err.Error())
		return err
	}

	err = bn.transporter.Subscribe(transporter.ChanMessageResponse, bn.ReceiveResponse)
	if err != nil {
		log.Warning(err.Error())
		return err
	}

	log.Default("# Name: " + string(bn.nodeVersionName) + " is running")
	log.Default("# NodeID: " + string(bn.nodeID))

	payload := transporter.PayloadDiscovery{
		Event:  transporter.EventConnected,
		Blocks: bn.Blocks(),
	}

	bn.sendDiscovery(&payload)

	go bn.daemon(bn.daemonChan)

	return nil
}

// AddBlock adds new Block struct to BlockNode blocks map
func (bn *BlockNode) AddBlock(blocks ...IBlock) error {
	var bk IBlock

	for _, b := range blocks {
		if bn.blocks[b.GetName()] != bk {
			return errors.ErrBlockAdded
		}

		bn.blocks[b.GetName()] = b
		log.Debug("block added: " + b.GetName().String())

		bn.trafficManager.AddDestination(types.NodeID(bn.nodeID), types.NodeVersionName(bn.nodeVersionName), b.GetName(), b.ActionsNames())
	}

	return nil
}

func (bn *BlockNode) Blocks() map[types.BlockName][]types.ActionName {
	var blocks map[types.BlockName][]types.ActionName = make(map[types.BlockName][]types.ActionName)

	for name, blk := range bn.blocks {
		actions := blk.Actions()

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

	bn.sendDiscovery(&payload)

	bn.daemonChan <- 1

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

func (bn *BlockNode) VersionName() types.NodeVersionName {
	return bn.nodeVersionName
}

func (bn *BlockNode) Database() *db.Database {
	return bn.database
}

func (bn *BlockNode) loadDatabase() {
	bn.database = db.NewDatabase(&bn.config.Database)
}

func (bn *BlockNode) daemon(daemonChan chan uint) {
	log.Debug("BlockNode daemon start")

	ticker := time.NewTicker(time.Duration(bn.config.Transporter.HeartbeatInterval) * time.Second)
L:
	for {
		select {
		case <-ticker.C:
			// send heartbeat
			payload := transporter.PayloadDiscovery{
				Event:  transporter.EventHeartbeat,
				Blocks: bn.Blocks(),
			}

			bn.sendDiscovery(&payload)

			for nodeID, lastSeen := range bn.trafficManager.Nodes() {
				if nodeID == bn.nodeID {
					continue
				}

				dateDiff := time.Since(lastSeen).Seconds()
				if dateDiff > float64(bn.config.Transporter.HeartbeatInterval+1) {
					bn.trafficManager.RemoveNodeID(nodeID)
				}
			}

		case <-daemonChan:
			break L
		}
	}

	log.Debug("BlockNode daemon quit")
}

func (bn *BlockNode) Receive(payload []byte) {
	log.Debug("Network.Receive []byte")

	pocket, err := DecodePocket(payload)
	if err != nil {
		log.Warning(err.Error())
		return
	}

	if pocket.FromID == bn.nodeID {
		return
	}

	if pocket.TargetID != nil && *pocket.TargetID != "" && *pocket.TargetID != bn.nodeID {
		log.Debug("not a target, skipping")
		return
	}

	if pocket.Channel == transporter.ChanDiscovery {
		if pocket.FromID == bn.nodeID {
			return
		}

		discoveryPayload, err := DecodePayload[transporter.PayloadDiscovery](pocket.Payload)
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

		bn.ProcessPocketDiscovery(newPocket)
		return
	}

	if pocket.Channel == transporter.ChanMessage {
		if bn.trafficManager.DestinationExist(*pocket.TargetAction) {
			messagePayload, err := DecodePayload[transporter.PayloadMessage](pocket.Payload)
			if err != nil {
				log.Warning(err.Error())
				return
			}

			response, err := bn.Send(&messagePayload, pocket.TargetAction)
			if err != nil {
				log.Warning(err.Error())
				return
			}

			responsePocket := transporter.NewPocket(transporter.ChanMessageResponse, bn.nodeVersionName, bn.nodeID, &pocket.FromID, nil, *response)
			responsePocket.ResponseHash = pocket.Hash

			responseJSON, err := json.Marshal(responsePocket)
			if err != nil {
				log.Warning(err.Error())
				return
			}

			err = bn.transporter.Send(responsePocket.Channel, responseJSON)
			if err != nil {
				log.Warning(err.Error())
				return
			}
		}
		return
	}
}

func (bn *BlockNode) ReceiveResponse(payload []byte) {
	log.Debug("Network.ReceiveResponse []byte")

	pocket, err := DecodePocket(payload)
	if err != nil {
		log.Warning(err.Error())
		return
	}

	if pocket.FromID == bn.nodeID {
		return
	}

	bn.sentHashesMutex.Lock()
	if bn.sentHashes[pocket.ResponseHash] != nil {
		bn.sentHashes[pocket.ResponseHash].responseChan <- pocket
	}
	bn.sentHashesMutex.Unlock()
}

// NewBlockNode creates new BlockNode struct
func NewBlockNode(options *BlockNodeOptions) *BlockNode {
	if blockNode == nil {
		lock.Lock()
		defer lock.Unlock()

		nodeID := common.CreateNodeID(options.Version, options.Name)

		blockNode = &BlockNode{
			nodeID:          nodeID,
			nodeVersionName: common.CreateNodeVersionName(options.Version, options.Name),
			options:         *options,
			blocks:          make(map[types.BlockName]IBlock),
			transporter:     nil,
			trafficManager:  traffic.NewManager(&nodeID),
			daemonChan:      make(chan uint),
			sentHashes:      map[string]*SentHash{},
			sentHashesMutex: new(sync.Mutex),
		}

		return blockNode
	}

	panic(errors.ErrBlockNodeInstantiated)
}
