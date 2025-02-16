package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgt4l/avito_shop/internal/avito_shop/auth"
	"github.com/dgt4l/avito_shop/internal/avito_shop/controller"
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	repository "github.com/dgt4l/avito_shop/internal/avito_shop/repository/pgsql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
	port        string
}

func NewShopHandler(srv ShopService, auth auth.AuthService, port string) *ShopHandler {
	e := echo.New()
	return &ShopHandler{
		e:           e,
		shopService: srv,
		auth:        auth,
		port:        port,
	}

}

func (h *ShopHandler) Start() error {
	RegisterRoutes(h)

	if err := h.e.Start(":" + h.port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.e.Logger.Fatal("Shutting down the server")
	}

	return nil
}

func (h *ShopHandler) Close(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}

func (h *ShopHandler) BuyItem(ctx echo.Context) error {
	const op = "internal.avito.shop.handler.ShopHandler.BuyItem"

	logrus.WithFields(logrus.Fields{"event": op}).Info(ctx.Get("id"))

	var request dto.BuyItemRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: ErrInvalidDataType.Error()})
	}

	var ok bool
	request.Id, ok = ctx.Get("id").(int)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	err := h.shopService.BuyItem(ctx.Request().Context(), &request)
	if err != nil && (errors.Is(err, repository.ErrNotEnoughCoins) ||
		errors.Is(err, repository.ErrItemNotFound) ||
		errors.Is(err, controller.ErrEmptyItemName)) {
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: err.Error()})
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{"event": op, "request": request}).Error(err)

		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *ShopHandler) SendCoin(ctx echo.Context) error {
	const op = "internal.avito.shop.handler.ShopHandler.SendCoin"

	logrus.WithFields(logrus.Fields{"event": op}).Info(ctx.Get("id"))

	var request dto.SendCoinRequest

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: ErrInvalidDataType.Error()})
	}

	fromUserId, ok := ctx.Get("id").(int)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	err := h.shopService.SendCoin(ctx.Request().Context(), fromUserId, &request)
	if err != nil && (errors.Is(err, controller.ErrShortUsername) ||
		errors.Is(err, controller.ErrInvalidAmount) ||
		errors.Is(err, repository.ErrNotEnoughCoins) ||
		errors.Is(err, repository.ErrUserToNotFound)) {
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: err.Error()})
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{"event": op, "fromUser": fromUserId, "request": request}).Error(err)

		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *ShopHandler) AuthUser(ctx echo.Context) error {
	const op = "internal.avito.shop.handler.ShopHandler.AuthUser"

	var request dto.AuthRequest

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: ErrInvalidDataType.Error()})
	}

	logrus.WithFields(logrus.Fields{"event": op}).Info(request)

	response, err := h.shopService.AuthUser(ctx.Request().Context(), &request)
	if err != nil && (errors.Is(err, controller.ErrShortPassword) ||
		errors.Is(err, controller.ErrShortUsername) ||
		errors.Is(err, controller.ErrInvalidPasswd)) {
		return ctx.JSON(http.StatusBadRequest, dto.BadRequestResponse{Errors: err.Error()})
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{"event": op, "request": request}).Error(err)

		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *ShopHandler) GetInfo(ctx echo.Context) error {
	const op = "internal.avito.shop.handler.ShopHandler.GetInfo"

	userId, ok := ctx.Get("id").(int)

	logrus.WithFields(logrus.Fields{"event": op}).Info(userId)

	if !ok {
		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	response, err := h.shopService.GetInfo(ctx.Request().Context(), userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{"event": op, "userId": userId}).Error(err)

		return ctx.JSON(http.StatusInternalServerError, dto.InternalServerErrorResponse{Errors: ErrInternalServer.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *ShopHandler) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
