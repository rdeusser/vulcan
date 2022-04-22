package scaffold

import "errors"

var (
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrUnknownAction     = errors.New("unknown action")
)
