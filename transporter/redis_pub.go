package transporter

import (
	"fmt"

	"github.com/block-api/block-node/log"
	"github.com/go-redis/redis/v8"
)

type PubDaemon struct {
	channel    Channel
	daemonChan chan DaemonCmd
	dataChan   chan interface{}
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
				str := fmt.Sprintf("%v", data)
				res := pd.redisPub.Publish(ctx, string(pd.channel), str)

				if err := res.Err(); err != nil {
					// return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Publish")
					log.Warning(err.Error())
				}
				fmt.Println(res)
			case cmd := <-pd.daemonChan:
				if cmd == DaemonStop {
					break B
				}
			}
		}

		defer close(pd.daemonChan)
		defer close(pd.dataChan)
		defer close(pd.outChan)
		defer pd.redisPub.Close()
	}(pd)
}

func NewPubDaemon(channel Channel, redisClient *redis.Client) PubDaemon {
	return PubDaemon{
		channel:    channel,
		daemonChan: make(chan DaemonCmd),
		dataChan:   make(chan interface{}),
		outChan:    make(chan<- interface{}),
		redisPub:   redisClient,
	}
}
