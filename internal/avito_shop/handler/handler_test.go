package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/dgt4l/avito_shop/test/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShopHandlerBuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	req := httptest.NewRequest(http.MethodGet, "/buy?item=test-item", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", 1)

	mockShopService.EXPECT().
		BuyItem(c.Request().Context(), &dto.BuyItemRequest{Id: 1, Item: "test-item"}).
		Return(nil)

	err := handler.BuyItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestShopHandlerSendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBufferString(`{"toUser":"user2","amount":50}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", 1)

	mockShopService.EXPECT().
		SendCoin(c.Request().Context(), 1, &dto.SendCoinRequest{ToUser: "user2", Amount: 50}).
		Return(nil)

	err := handler.SendCoin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestShopHandlerAuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	requestBody := `{"username":"testuser","password":"testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockShopService.EXPECT().
		AuthUser(c.Request().Context(), &dto.AuthRequest{Username: "testuser", Password: "testpassword"}).
		Return(&dto.AuthResponse{Token: "test-token"}, nil)

	err := handler.AuthUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"token":"test-token"`)
}

func TestShopHandlerGetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", 1)

	expectedResponse := &dto.InfoResponse{
		Coins: 100,
		Inventory: []dto.Inventory{
			{Type: "item1", Quantity: 1},
		},
		CoinHistory: dto.CoinHistory{
			Received: []dto.Received{
				{FromUser: "user2", Amount: 50},
			},
			Sent: []dto.Sent{
				{ToUser: "user3", Amount: 30},
			},
		},
	}
	mockShopService.EXPECT().
		GetInfo(c.Request().Context(), 1).
		Return(expectedResponse, nil)

	err := handler.GetInfo(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"coins":100`)
}

func TestShopHandlerPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Ping(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

func TestShopHandlerAuthMiddleware_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	e.Use(handler.AuthMiddleware())

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	mockAuthService.EXPECT().
		ParseToken("valid-token").
		Return(1, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(authorizationHeader, "Bearer valid-token")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "success", rec.Body.String())
}

func TestShopHandlerAuthMiddleware_EmptyToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	e.Use(handler.AuthMiddleware())

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrEmptyToken.Error())
}

func TestShopHandlerAuthMiddleware_InvalidAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	e.Use(handler.AuthMiddleware())

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(authorizationHeader, "InvalidHeader")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInvalidAuthHeader.Error())
}

func TestShopHandlerAuthMiddleware_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShopService := mocks.NewMockShopService(ctrl)
	mockAuthService := mocks.NewMockAuthService(ctrl)

	e := echo.New()
	handler := NewShopHandler(mockShopService, mockAuthService, "8080")

	e.Use(handler.AuthMiddleware())

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	mockAuthService.EXPECT().
		ParseToken("invalid-token").
		Return(0, errors.New("invalid token"))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(authorizationHeader, "Bearer invalid-token")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), ErrInternalServer.Error())
}
