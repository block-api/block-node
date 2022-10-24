package transporter

import (
	"encoding/json"

	"github.com/block-api/block-node/log"
	"github.com/go-redis/redis/v8"
)

func NewSubDaemon(nodeID string, channel Channel, redisSub *redis.PubSub) SubDaemon {
	return SubDaemon{
		nodeID:     nodeID,
		channel:    channel,
		daemonChan: make(chan DaemonCmd),
		outChan:    make(chan Pocket[[]byte]),
		redisSub:   redisSub,
	}
}

func (sd *SubDaemon) Start(callback func(pocket Pocket[[]byte])) {
	go func(sd *SubDaemon) {
		log.Debug("SubDaemon started: " + string(sd.channel))

		// B:
		for msg := range sd.redisSub.Channel() {
			switch msg.Channel {
			case string(sd.channel):
				pocket, err := sd.processPayload(msg.Payload)
				if err != nil {
					log.Warning(err.Error())
					break
				}
				log.Debug("parsed new transporter Pocket")

				if (pocket.FromID == sd.nodeID) || (pocket.TargetID != "" && pocket.TargetID != sd.nodeID) {
					log.Debug("skip self/not target")
					break
				}

				sd.outChan <- *pocket
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

func (sd *SubDaemon) processPayload(payload string) (*Pocket[[]byte], error) {
	pocket := Pocket[[]byte]{}
	err := json.Unmarshal([]byte(payload), &pocket)
	if err != nil {
		return nil, err
	}

	return &pocket, nil
}

func (sd *SubDaemon) OutChat() chan Pocket[[]byte] {
	return sd.outChan
}
