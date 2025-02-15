package repository

import "errors"

var ErrUserNotFound = errors.New("user not found")

var ErrNotEnoughCoins = errors.New("not enough coins")

var ErrItemNotFound = errors.New("item not found")

var ErrUserToNotFound = errors.New("user receiver not found")
