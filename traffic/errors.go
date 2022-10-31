package traffic

import "errors"

var (
	ErrNoNodeFound = errors.New("can't find any available nodes for target action")
)
