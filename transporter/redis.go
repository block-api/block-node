package transporter

import (
	"context"
	"fmt"

	"github.com/block-api/block-node/config"
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

func (r *Redis) ConneDisconnectct() error {
	return nil
}

func (r *Redis) Send(target string, payload interface{}) error {
	return nil
}

// Disconnect() error
// Send(target string, payload interface{}) error

func NewRedis(configRedis *config.ConfigRedisTransporter) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configRedis.Host + ":" + configRedis.Port,
		Username: configRedis.Username,
		Password: configRedis.Password, // no password set
		DB:       int(configRedis.Db),  // use default DB
	})

	fmt.Println(rdb)
	return &Redis{}, nil
}
