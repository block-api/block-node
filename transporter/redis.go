package transporter

import (
	"context"

	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/log"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Redis struct {
	Transporter
	subClient *redis.Client
	pubClient *redis.Client
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

func (r *Redis) Send(target string, payload interface{}) error {
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

	return &Redis{
		subClient: redisSub,
		pubClient: redisPub,
	}, nil
}
