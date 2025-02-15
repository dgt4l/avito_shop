package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgt4l/avito_shop/internal/avito_shop/auth"
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/labstack/echo/v4"
)

type ShopService interface {
	GetInfo(ctx context.Context, userId int) (*dto.InfoResponse, error)
	BuyItem(ctx context.Context, request *dto.BuyItemRequest) error
	AuthUser(ctx context.Context, request *dto.AuthRequest) (*dto.AuthResponse, error)
	SendCoin(ctx context.Context, fromUserId int, request *dto.SendCoinRequest) error
}

type ShopHandler struct {
	e           *echo.Echo
	shopService ShopService
	auth        auth.AuthService
}

func NewShopHandler(srv ShopService, auth auth.AuthService) *ShopHandler {
	e := echo.New()
	return &ShopHandler{
		e:           e,
		shopService: srv,
		auth:        auth,
	}

}

func (h *ShopHandler) Start() error {
	RegisterRoutes(h)

	if err := h.e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.e.Logger.Fatal("Shutting down the server")
	}

	return nil
}

func (h *ShopHandler) Close(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}

func (h *ShopHandler) BuyItem(ctx echo.Context) error {
	var request dto.BuyItemRequest
	request.Item = ctx.QueryParam("item")
	request.Id = ctx.Get("id").(int)
	err := h.shopService.BuyItem(ctx.Request().Context(), &request)
	if err != nil {
		ctx.Set("errors", err.Error())
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: err.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *ShopHandler) SendCoin(ctx echo.Context) error {
	var request dto.SendCoinRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Set("errors", err.Error())
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: err.Error()})
	}
	fromUserId := ctx.Get("id").(int)
	err := h.shopService.SendCoin(ctx.Request().Context(), fromUserId, &request)
	if err != nil {
		ctx.Set("errors", err.Error())
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: err.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *ShopHandler) AuthUser(ctx echo.Context) error {
	var request dto.AuthRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Set("errors", err.Error())
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: err.Error()})
	}

	response, err := h.shopService.AuthUser(ctx.Request().Context(), &request)
	if err != nil {
		ctx.Set("errors", err.Error())
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *ShopHandler) GetInfo(ctx echo.Context) error {
	userId := ctx.Get("id").(int)
	response, err := h.shopService.GetInfo(ctx.Request().Context(), userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *ShopHandler) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
