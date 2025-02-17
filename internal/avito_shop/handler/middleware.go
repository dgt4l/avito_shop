package handler

import (
	"net/http"
	"strings"

	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
)

func (h *ShopHandler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			const op = "internal.avito_shop.handler.AuthMiddleware"

			header := ctx.Request().Header.Get(authorizationHeader)
			if header == "" {
				return ctx.JSON(http.StatusUnauthorized, dto.UnauthorizedResponse{Errors: ErrEmptyToken.Error()})
			}

			headerSplit := strings.Split(header, " ")
			if len(headerSplit) != 2 {
				return ctx.JSON(http.StatusUnauthorized, dto.UnauthorizedResponse{Errors: ErrInvalidAuthHeader.Error()})
			}

			id, err := h.auth.ParseToken(headerSplit[1])
			if err != nil {
				logrus.WithFields(logrus.Fields{"event": op}).Error(err)

				return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
			}

			ctx.Set("id", id)
			return next(ctx)
		}
	}
}
