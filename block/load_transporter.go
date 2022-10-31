package block

import (
	"github.com/block-api/block-node/config"
	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/transporter"
)

// loadTransporter is responsible for loading one of configured & supported transporters from config
func (bn *BlockNode) loadTransporter(transporterConfig config.Transporter) error {
	if transporterConfig.Redis != nil {
		redisTransport, err := transporter.NewRedis(bn.nodeID, transporterConfig.Redis)
		if err != nil {
			panic(err)
		}

		bn.transporter = redisTransport
		return nil
	}

	return errors.ErrTransporterNotConfigured
}
