package handler

import "errors"

var InternalServerError = errors.New("internal server error")

var ErrEmptyItemName = errors.New("empty item name")

var ErrEmptyUsername = errors.New("username is empty")

var ErrShortPassword = errors.New("password is too short")

var ErrInvalidAmount = errors.New("amount must be positive number")
