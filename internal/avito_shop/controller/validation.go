package controller

import (
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
)

func ValidateAuth(request *dto.AuthRequest) error {
	if len(request.Username) < 4 {
		return ErrShortUsername
	}

	if len(request.Password) < 8 {
		return ErrShortPassword
	}

	return nil
}

func ValidateSendCoin(request *dto.SendCoinRequest) error {
	if len(request.ToUser) < 4 {
		return ErrShortUsername
	}

	if request.Amount <= 0 {
		return ErrInvalidAmount
	}

	return nil
}

func ValidateBuyItem(request *dto.BuyItemRequest) error {
	if request.Item == "" {
		return ErrEmptyItemName
	}

	return nil
}
