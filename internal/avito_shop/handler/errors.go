package handler

import (
	"errors"
)

var ErrInternalServer = errors.New("internal error")

var ErrEmptyItemName = errors.New("empty item name")

var ErrInvalidDataType = errors.New("invalid type of data")
