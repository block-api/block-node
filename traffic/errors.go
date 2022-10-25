package traffic

import "errors"

var (
	ErrNodeNameIDExist    = errors.New("node id for given name already exist")
	ErrNodeNameIDNotExist = errors.New("node id for given name does not exist")
)
