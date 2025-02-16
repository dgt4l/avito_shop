package handler

import (
	"errors"
)

var ErrInternalServer = errors.New("internal error")

var ErrInvalidDataType = errors.New("invalid type of data")

var ErrEmptyToken = errors.New("empty token")

var ErrInvalidAuthHeader = errors.New("invalid auth header")

var ErrInvalidToken = errors.New("invalid token")
