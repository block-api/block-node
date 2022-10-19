package block

import (
	"fmt"

	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/transporter"
)

func (bn *BlockNode) loadTransporter(transporterConfig config.ConfigTransporter) error {

	if transporterConfig.Redis != nil {
		redisTransport, err := transporter.NewRedis(transporterConfig.Redis)
		if err != nil {
			panic(err)
		}

		fmt.Println(redisTransport)
		bn.transporter = redisTransport
	}

	return nil
}
