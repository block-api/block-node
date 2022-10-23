package transporter

import (
	"github.com/block-api/block-node/log"
	"github.com/go-redis/redis/v8"
)

type PubDaemon struct {
	channel    Channel
	daemonChan chan DaemonCmd
	dataChan   chan []byte
	outChan    chan<- interface{}
	redisPub   *redis.Client
}

func (pd *PubDaemon) Stop() {
	pd.daemonChan <- DaemonStop
}

func (pd *PubDaemon) Start() {
	go func(pd *PubDaemon) {
		log.Debug("PubDaemon started " + string(pd.channel))
	B:
		for {
			select {
			case data := <-pd.dataChan:
				res := pd.redisPub.Publish(ctx, string(pd.channel), data)
				if err := res.Err(); err != nil {
					log.Warning(err.Error())
				}
			case cmd := <-pd.daemonChan:
				if cmd == DaemonStop {
					break B
				}
			}
		}

		defer pd.redisPub.Close()
		defer close(pd.daemonChan)
		defer close(pd.dataChan)
		defer close(pd.outChan)
	}(pd)
}

func NewPubDaemon(channel Channel, redisClient *redis.Client) PubDaemon {
	return PubDaemon{
		channel:    channel,
		daemonChan: make(chan DaemonCmd),
		dataChan:   make(chan []byte),
		outChan:    make(chan<- interface{}),
		redisPub:   redisClient,
	}
}
