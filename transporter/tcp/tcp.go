package tcp

import "github.com/block-api/block-node/transporter"

type Tcp struct {
	transporter.Transporter
}

func New() *Tcp {
	tcp := &Tcp{}

	return tcp
}
