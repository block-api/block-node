package transporter

import (
	"context"

	"github.com/block-api/block-node/common/types"
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
	nodeID     types.NodeID
	channel    Channel
	daemonChan chan DaemonCmd
	outChan    chan []byte
	redisSub   *redis.PubSub
}

type Redis struct {
	Transporter
	nodeID    types.NodeID
	subClient *redis.Client
	pubClient *redis.Client
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

func (r *Redis) Send(channel Channel, payload []byte) error {
	if r.pubDaemon[string(channel)].dataChan == nil {
		pubDaemon := NewPubDaemon(channel, r.pubClient)
		pubDaemon.Start()

		r.pubDaemon[string(channel)] = pubDaemon
		r.pubDaemon[string(channel)].dataChan <- payload

		log.Debug("published to: " + string(channel))

		return nil
	}

	log.Debug("published to: " + string(channel))

	r.pubDaemon[string(channel)].dataChan <- payload

	return nil
}

func (r *Redis) Subscribe(channel Channel, callback func(payload []byte)) error {
	sub := r.subClient.Subscribe(ctx, string(channel))

	subDaemon := NewSubDaemon(r.nodeID, channel, sub)
	subDaemon.Start(callback)

	r.subDaemon = append(r.subDaemon, subDaemon)

	return nil
}

func NewRedis(nodeID types.NodeID, configRedis *config.RedisTransporter) (*Redis, error) {
	redisOptions := &redis.Options{
		Addr:     configRedis.Host + ":" + configRedis.Port,
		Username: configRedis.Username,
		Password: configRedis.Password,
		DB:       int(configRedis.Db),
	}

	redisSub := redis.NewClient(redisOptions)
	redisPub := redis.NewClient(redisOptions)

	log.Debug("redis instances created: " + configRedis.Host + ":" + configRedis.Port)

	return &Redis{
		nodeID:    nodeID,
		subClient: redisSub,
		pubClient: redisPub,
		pubDaemon: make(map[string]PubDaemon),
	}, nil
}
