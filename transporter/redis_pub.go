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
		// var err error
	B:
		for {
			select {
			case data := <-pd.dataChan:
				// fmt.Println(data)
				// var bData bytes.Buffer

				// if err := json.NewEncoder(&bData).Encode(err); err != nil {
				// 	log.Warning(err.Error())
				// 	break
				// }
				// log.Warning(bData.String())
				// fmt.Println(bData.String())
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
			// s
			// switch msg.Channel {
			// case string(sd.channel):
			// 	fmt.Println(msg.Payload)
			// }
		}
		// for {
		// 	select {
		// 	case cmd := <-sd.daemonChan:
		// 		fmt.Println(cmd)
		// 		if cmd == DaemonStop {
		// 			break B
		// 		}
		// 	case m := msg:
		// 		if err != nil {
		// 			log.Debug(err)
		// 			break
		// 		}
		// 		// case data := <-sd.dataChan:
		// 		// 	fmt.Println(data)
		// 	}
		// }

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
