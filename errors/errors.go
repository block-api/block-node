package errors

import (
	"errors"
)

const (
	ErrNotAuthorized         = "not authorized"
	ErrNotFound              = "not found"
	ErrBlockAdded            = "block has been already added"
	ErrBlockNodeInstantiated = "BlockNode instantiated"
)

func New(text string) error {
	return errors.New(text)
}
