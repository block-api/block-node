package errors

import (
	"errors"
)

const (
	ErrNotAuthorized            = "not authorized"
	ErrNotFound                 = "not found"
	ErrBlockAdded               = "block has been already added"
	ErrBlockNodeInstantiated    = "BlockNode instantiated"
	ErrTransporterNotConfigured = "transporter is not configured in config.yml file"
)

var ErrInvalidTargetAction = errors.New("invalid target action")

func New(text string) error {
	return errors.New(text)
}
