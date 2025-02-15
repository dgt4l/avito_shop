package handler

import (
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
)

func ValidateAuth(request *dto.AuthRequest) error {
	if request.Username == "" {
		return ErrEmptyUsername
	}

	if len(request.Password) < 8 {
		return ErrShortPassword
	}

	return nil
}

func ValidateSendCoin(request *dto.SendCoinRequest) error {
	if request.ToUser == "" {
		return ErrEmptyUsername
	}

	if request.Amount <= 0 {
		return ErrInvalidAmount
	}

	return nil
}
