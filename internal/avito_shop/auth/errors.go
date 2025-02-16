package auth

import "errors"

var ErrInvalidSignMethod = errors.New("invalid signing method")

var ErrClaimIdFails = errors.New("claim parsing id fails")

var ErrClaimMissing = errors.New("claim missing")
