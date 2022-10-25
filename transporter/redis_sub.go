package transporter

import (
	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/utils"
	"github.com/go-redis/redis/v8"
)

func NewSubDaemon(nodeID utils.NodeID, channel Channel, redisSub *redis.PubSub) SubDaemon {
	return SubDaemon{
		nodeID:     nodeID,
		channel:    channel,
		daemonChan: make(chan DaemonCmd),
		outChan:    make(chan []byte),
		redisSub:   redisSub,
	}
}

func (sd *SubDaemon) Start(callback func(payload []byte)) {
	go func(sd *SubDaemon) {
		log.Debug("SubDaemon started: " + string(sd.channel))

		// B:
		for msg := range sd.redisSub.Channel() {
			switch msg.Channel {
			case string(sd.channel):
				sd.outChan <- []byte(msg.Payload)
			}
		}

		defer close(sd.daemonChan)
		defer close(sd.outChan)
		defer sd.redisSub.Close()
	}(sd)

	go func(sd *SubDaemon) {
		for pocket := range sd.outChan {
			callback(pocket)
		}
	}(sd)
}

func (sd *SubDaemon) Stop() {
	sd.daemonChan <- DaemonStop
}

func (sd *SubDaemon) OutChan() chan []byte {
	return sd.outChan
}
