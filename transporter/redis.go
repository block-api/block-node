package transporter

import (
	"context"
	"fmt"

	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/log"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type DaemonCmd uint

const (
	DaemonStop DaemonCmd = iota
	DaemonStart
)

type SubDaemon struct {
	channel    Channel
	daemonChan chan DaemonCmd
	// dataChan   chan interface{}
	outChan  chan<- interface{}
	redisSub *redis.PubSub
}

type Redis struct {
	Transporter
	subClient *redis.Client
	pubClient *redis.Client
	// subClients []*redis.PubSub
	subDaemon []SubDaemon
	pubDaemon map[string]PubDaemon
}

func (r *Redis) Connect() error {
	return nil
}

func (r *Redis) Disconnect() error {
	log.Debug("closing redis connections")

	err := r.subClient.Close()
	if err != nil {
		return err
	}

	return r.pubClient.Close()
}

func (r *Redis) Send(channel Channel, payload interface{}) error {
	fmt.Println(r.pubDaemon[string(channel)])

	if r.pubDaemon[string(channel)].dataChan == nil {
		pubDaemon := NewPubDaemon(channel, r.pubClient)
		pubDaemon.Start()

		pubDaemon.dataChan <- payload
		return nil
	}

	r.pubDaemon[string(channel)].dataChan <- payload

	return nil
}

func (r *Redis) Subscribe(channel Channel) error {
	sub := r.subClient.Subscribe(ctx, string(channel))

	subDaemon := NewSubDaemon(channel, sub)
	subDaemon.Start()

	r.subDaemon = append(r.subDaemon, subDaemon)

	return nil
}

func NewRedis(configRedis *config.ConfigRedisTransporter) (*Redis, error) {
	redisOptions := &redis.Options{
		Addr:     configRedis.Host + ":" + configRedis.Port,
		Username: configRedis.Username,
		Password: configRedis.Password, // no password set
		DB:       int(configRedis.Db),  // use default DB
	}

	redisSub := redis.NewClient(redisOptions)
	redisPub := redis.NewClient(redisOptions)

	log.Debug("redis instances created: " + configRedis.Host + ":" + configRedis.Port)

	return &Redis{
		subClient: redisSub,
		pubClient: redisPub,
		pubDaemon: make(map[string]PubDaemon),
	}, nil
}

// -- SubDaemon -- //

func (sd *SubDaemon) Start() {
	go func(sd *SubDaemon) {
		log.Debug("SubDaemon started: " + string(sd.channel))

		// B:
		for msg := range sd.redisSub.Channel() {
			switch msg.Channel {
			case string(sd.channel):
				log.Warning("---> receive")
				log.Warning(msg.Payload)
			}
		}

		defer close(sd.daemonChan)
		// defer close(sd.dataChan)
		defer close(sd.outChan)
		defer sd.redisSub.Close()
	}(sd)
}

func (sd *SubDaemon) Stop() {
	sd.daemonChan <- DaemonStop
}

func NewSubDaemon(channel Channel, redisSub *redis.PubSub) SubDaemon {
	return SubDaemon{
		channel:    channel,
		daemonChan: make(chan DaemonCmd),
		// dataChan:   make(chan interface{}),
		outChan:  make(chan<- interface{}),
		redisSub: redisSub,
	}
}
