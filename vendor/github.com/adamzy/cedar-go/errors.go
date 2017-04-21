package cedar

import "errors"

var (
	ErrInvalidDataType = errors.New("cedar: invalid datatype")
	ErrInvalidValue    = errors.New("cedar: invalid value")
	ErrInvalidKey      = errors.New("cedar: invalid key")
	ErrNoPath          = errors.New("cedar: no path")
	ErrNoValue         = errors.New("cedar: no value")
)
