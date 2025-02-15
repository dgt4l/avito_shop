package handler

import (
	"net/http"
	"strings"

	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/labstack/echo/v4"
)

const (
	authorizationHeader = "Authorization"
)

func (h *ShopHandler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header.Get(authorizationHeader)
			if header == "" {
				return ctx.JSON(http.StatusUnauthorized, dto.UnauthorizedResponse{Errors: "empty token"})
			}

			headerSplit := strings.Split(header, " ")
			if len(headerSplit) != 2 {
				return ctx.JSON(http.StatusUnauthorized, dto.UnauthorizedResponse{Errors: "invalid auth header"})
			}

			id, err := h.auth.ParseToken(headerSplit[1])
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, dto.UnauthorizedResponse{Errors: err.Error()})
			}

			ctx.Set("id", id)
			return next(ctx)
		}
	}
}
