package errors

import (
	"errors"
)

var (
	ErrInvalidTargetAction      = errors.New("invalid target action")
	ErrBlockAdded               = errors.New("block has been already added")
	ErrBlockNodeInstantiated    = errors.New("BlockNode instantiated")
	ErrTransporterNotConfigured = errors.New("transporter is not configured in config.yml file")
	ErrInvalidFileType          = errors.New("invalid file type")
	ErrConfigDirMissing         = errors.New("CONFIG_DIR is missing")
)
