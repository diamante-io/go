package sequence

import (
	"errors"
)

var (
	ErrNoMoreRoom  = errors.New("queue full")
	ErrBadSequence = errors.New("bad sequence")
)
