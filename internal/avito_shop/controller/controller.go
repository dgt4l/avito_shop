package controller

import (
	"context"

	"github.com/dgt4l/avito_shop/internal/avito_shop/auth"
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/dgt4l/avito_shop/internal/avito_shop/models"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	BuyItem(ctx context.Context, id int, item string) error
	GetInfo(ctx context.Context, userId int) (*dto.InfoResponse, error)
	SendCoin(ctx context.Context, toUser string, fromUserId, amount int) error
	CreateUser(ctx context.Context, username, password string) (int, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
}

type ShopService struct {
	repo Repository
	auth auth.AuthService
	cfg  ServiceConfig
}

func NewShopService(repo Repository, auth auth.AuthService, cfg ServiceConfig) *ShopService {
	return &ShopService{
		repo: repo,
		auth: auth,
		cfg:  cfg,
	}
}

func (s *ShopService) BuyItem(ctx context.Context, request *dto.BuyItemRequest) error {
	return s.repo.BuyItem(ctx, request.Id, request.Item)
}

func (s *ShopService) GetInfo(ctx context.Context, userId int) (*dto.InfoResponse, error) {
	return s.repo.GetInfo(ctx, userId)
}

func (s *ShopService) SendCoin(ctx context.Context, fromUserId int, request *dto.SendCoinRequest) error {
	if err := ValidateSendCoin(request); err != nil {
		return err
	}
	return s.repo.SendCoin(ctx, request.ToUser, fromUserId, request.Amount)
}

func (s *ShopService) CreateUser(ctx context.Context, request *dto.AuthRequest) (*models.User, error) {
	password_hash, err := s.generatePasswordHash(request.Password)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.CreateUser(ctx, request.Username, password_hash)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Id:       id,
		Username: request.Username,
		Password: request.Password,
	}

	return &user, nil
}

func (s *ShopService) generatePasswordHash(password string) (string, error) {
	var passwordBytes = []byte(password + s.cfg.Salt)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

func (s *ShopService) AuthUser(ctx context.Context, request *dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := ValidateAuth(request); err != nil {
		return nil, err
	}
	user, err := s.repo.GetUser(ctx, request.Username)
	if err != nil {
		user, err = s.CreateUser(ctx, request)
		if err != nil {
			return nil, err
		}

		token, err := s.auth.GenerateToken(user)
		if err != nil {
			return nil, err
		}
		return &dto.AuthResponse{Token: token}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password+s.cfg.Salt)); err != nil {
		return nil, ErrInvalidPasswd
	}

	token, err := s.auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{Token: token}, nil
}
