package controller

import (
	"context"
	"testing"

	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/dgt4l/avito_shop/internal/avito_shop/models"
	repository "github.com/dgt4l/avito_shop/internal/avito_shop/repository/pgsql"
	"github.com/dgt4l/avito_shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestShopService_BuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	req := &dto.BuyItemRequest{Id: 1, Item: "item1"}

	mockRepo.EXPECT().BuyItem(ctx, req.Id, req.Item).Return(nil)

	err := service.BuyItem(ctx, req)
	assert.NoError(t, err)
}

func TestShopService_GetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	userId := 1
	expectedInfo := &dto.InfoResponse{
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

	mockRepo.EXPECT().GetInfo(ctx, userId).Return(expectedInfo, nil)

	info, err := service.GetInfo(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedInfo, info)
}

func TestShopService_SendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	fromUserId := 1
	req := &dto.SendCoinRequest{ToUser: "user2", Amount: 50}

	mockRepo.EXPECT().SendCoin(ctx, req.ToUser, fromUserId, req.Amount).Return(nil)

	err := service.SendCoin(ctx, fromUserId, req)
	assert.NoError(t, err)
}

func TestShopService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	req := &dto.AuthRequest{Username: "user1", Password: "password1"}
	expectedUser := &models.User{Id: 1, Username: "user1", Password: "password1"}

	mockRepo.EXPECT().CreateUser(ctx, req.Username, gomock.Any()).Return(1, nil)

	user, err := service.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
}

func TestShopService_AuthUser_NewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	req := &dto.AuthRequest{Username: "user1", Password: "password1"}
	expectedUser := &models.User{Id: 1, Username: "user1", Password: "password1"}
	expectedToken := "test-token"

	mockRepo.EXPECT().GetUser(ctx, req.Username).Return(nil, repository.ErrUserNotFound)
	mockRepo.EXPECT().CreateUser(ctx, req.Username, gomock.Any()).Return(1, nil)
	mockAuth.EXPECT().GenerateToken(expectedUser).Return(expectedToken, nil)

	authResponse, err := service.AuthUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, authResponse.Token)
}

func TestShopService_AuthUser_ExistingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	req := &dto.AuthRequest{Username: "user1", Password: "password1"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1test-salt"), bcrypt.DefaultCost)
	existingUser := &models.User{Id: 1, Username: "user1", Password: string(hashedPassword)}
	expectedToken := "test-token"

	mockRepo.EXPECT().GetUser(ctx, req.Username).Return(existingUser, nil)
	mockAuth.EXPECT().GenerateToken(existingUser).Return(expectedToken, nil)

	authResponse, err := service.AuthUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, authResponse.Token)
}

func TestShopService_AuthUser_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)

	service := NewShopService(mockRepo, mockAuth, ServiceConfig{Salt: "test-salt"})

	ctx := context.Background()
	req := &dto.AuthRequest{Username: "user1", Password: "wrong-password"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1test-salt"), bcrypt.DefaultCost)
	existingUser := &models.User{Id: 1, Username: "user1", Password: string(hashedPassword)}

	mockRepo.EXPECT().GetUser(ctx, req.Username).Return(existingUser, nil)

	_, err := service.AuthUser(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPasswd, err)
}
