package storage

import "github.com/block-api/block-node/transporter"

type Storage struct {
	transporter *transporter.Transporter
}

func NewStorage(transporter *transporter.Transporter) (Storage, error) {

	return Storage{
		transporter,
	}, nil
}
