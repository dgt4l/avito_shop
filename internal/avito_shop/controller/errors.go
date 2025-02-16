package controller

import "errors"

var ErrInvalidPasswd = errors.New("invalid password")

var ErrShortPassword = errors.New("password is too short")

var ErrInvalidAmount = errors.New("amount must be positive number")

var ErrShortUsername = errors.New("username is too short")

var ErrEmptyItemName = errors.New("empty item name")
