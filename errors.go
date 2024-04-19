package console

import "errors"

var (
	ErrProcessNotStarted = errors.New("process has not been started")
	ErrInvalidCmd        = errors.New("invalid command")
)
