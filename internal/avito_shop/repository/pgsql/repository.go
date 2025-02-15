package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/dgt4l/avito_shop/internal/avito_shop/models"

	_ "github.com/lib/pq"
)

type Repository struct {
	db  *sqlx.DB
	cfg DBConfig
}

func NewRepository(config DBConfig) (*Repository, error) {
	db, err := sqlx.Open(
		"postgres", fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=%s",
			config.DBDriver,
			config.DBUser,
			config.DBPass,
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.DBSSL,
		),
	)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Repository{
		db:  db,
		cfg: config,
	}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) GetUser(ctx context.Context, username string) (*models.User, error) {
	logrus.Info("GetUser: ", username)
	var user models.User
	err := r.db.QueryRowxContext(ctx, getFromUsers, username).StructScan(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (r *Repository) BuyItem(ctx context.Context, id int, item string) error {
	logrus.Info("BuyItem: ", item)
	tx, err := r.db.BeginTxx(ctx, nil)
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			logrus.Error(err)
		}
	}()

	if err != nil {
		return err
	}

	var itemModel models.Item
	err = tx.QueryRowxContext(ctx, getFromItems, item).StructScan(&itemModel)
	if err != nil {
		return err
	}

	var userCoins int
	err = tx.QueryRowxContext(ctx, getCoinsFromUser, id).Scan(&userCoins)
	if err != nil {
		return err
	}

	if userCoins < itemModel.Price {
		return errors.New("Not enough coins")
	}

	_, err = tx.ExecContext(ctx, updateCoinsFromUser, itemModel.Price, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, insertToInventory, id, itemModel.Id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetInfo(ctx context.Context, userId int) (*dto.InfoResponse, error) {
	logrus.Info("GetInfo: ", userId)

	var userCoins int
	err := r.db.QueryRowxContext(ctx, getCoins, userId).Scan(&userCoins)
	if err != nil {
		return &dto.InfoResponse{}, err
	}

	var userInventory []dto.Inventory
	err = r.db.SelectContext(ctx, &userInventory, getUserInventory, userId)
	if err != nil {
		return &dto.InfoResponse{}, err
	}

	var userRecieved []dto.Received
	err = r.db.SelectContext(ctx, &userRecieved, getUserRecieved, userId)
	if err != nil {
		return &dto.InfoResponse{}, err
	}

	var userSent []dto.Sent
	err = r.db.SelectContext(ctx, &userSent, getUserSent, userId)
	if err != nil {
		return &dto.InfoResponse{}, err
	}

	return &dto.InfoResponse{
		Coins:     userCoins,
		Inventory: userInventory,
		CoinHistory: dto.CoinHistory{
			Received: userRecieved,
			Sent:     userSent,
		},
	}, nil
}

func (r *Repository) SendCoin(ctx context.Context, toUser string, fromUserId, amount int) error {
	logrus.Info("SendCoin: ", toUser, ", amount: ", amount)
	tx, err := r.db.BeginTxx(ctx, nil)

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			logrus.Error(err)
		}
	}()

	if err != nil {
		return err
	}

	var toUserId int
	err = tx.QueryRowContext(ctx, getIdFromUsers, toUser).Scan(&toUserId)

	var userCoins int
	err = tx.QueryRowxContext(ctx, getCoinsFromUser, fromUserId).Scan(&userCoins)
	logrus.Info("select coins fromUser")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, getCoinsToUser, toUser)
	logrus.Info("select coins toUser")
	if err != nil {
		return err
	}

	if userCoins < amount {
		return errors.New("Not enough coins")
	}

	_, err = tx.ExecContext(ctx, updateCoinsFromUser, amount, fromUserId)
	logrus.Info("update coins fromUser")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateCoinsToUser, amount, toUser)
	logrus.Info("update coins toUser")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, insertToTransactions, fromUserId, toUserId, amount)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateUser(ctx context.Context, username, password string) (int, error) {
	var id int
	err := r.db.QueryRowxContext(ctx, insertToUsers, username, password, r.cfg.DefaultCoins).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
