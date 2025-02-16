package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	"github.com/dgt4l/avito_shop/internal/avito_shop/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: sqlx.NewDb(db, "sqlmock")}

	tests := []struct {
		name         string
		username     string
		mockExpect   func()
		expectedResp func(*testing.T, *models.User, error)
	}{
		{
			name:     "success GetUser",
			username: "testuser",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password_salt"}).
					AddRow(1, "testuser", "hashedpassword")
				mock.ExpectQuery(regexp.QuoteMeta(getFromUsers)).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, user *models.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "hashedpassword", user.Password)
			},
		},
		{
			name:     "user not found",
			username: "nonexistent",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(getFromUsers)).
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, user *models.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrUserNotFound, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			user, err := repo.GetUser(context.Background(), tt.username)
			tt.expectedResp(t, user, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_BuyItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: sqlx.NewDb(db, "sqlmock")}

	tests := []struct {
		name         string
		userId       int
		item         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name:   "successful BuyItem",
			userId: 1,
			item:   "item1",
			mockExpect: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow(1, "item1", 100)
				mock.ExpectQuery(regexp.QuoteMeta(getFromItems)).
					WithArgs("item1").
					WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(getCoinsFromUser)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(200))
				mock.ExpectExec(regexp.QuoteMeta(updateCoinsFromUser)).
					WithArgs(100, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertToInventory)).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "item not found",
			userId: 1,
			item:   "nonexistent",
			mockExpect: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(getFromItems)).
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrItemNotFound, err)
			},
		},
		{
			name:   "not enough coins",
			userId: 1,
			item:   "item1",
			mockExpect: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow(1, "item1", 100)
				mock.ExpectQuery(regexp.QuoteMeta(getFromItems)).
					WithArgs("item1").
					WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(getCoinsFromUser)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(50))
				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrNotEnoughCoins, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			err := repo.BuyItem(context.Background(), tt.userId, tt.item)
			tt.expectedResp(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: sqlx.NewDb(db, "sqlmock")}

	tests := []struct {
		name         string
		userId       int
		mockExpect   func()
		expectedResp func(*testing.T, *dto.InfoResponse, error)
	}{
		{
			name:   "success GetInfo",
			userId: 1,
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(getCoins)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(100))
				mock.ExpectQuery(regexp.QuoteMeta(getUserInventory)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"name", "quantity"}).AddRow("item1", 1))
				mock.ExpectQuery(regexp.QuoteMeta(getUserRecieved)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"from_user", "amount"}).AddRow("user2", 50))
				mock.ExpectQuery(regexp.QuoteMeta(getUserSent)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user", "amount"}).AddRow("user3", 30))
			},
			expectedResp: func(t *testing.T, info *dto.InfoResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 100, info.Coins)
				assert.Equal(t, 1, len(info.Inventory))
				assert.Equal(t, 1, len(info.CoinHistory.Received))
				assert.Equal(t, 1, len(info.CoinHistory.Sent))
			},
		},
		{
			name:   "user not found",
			userId: 1,
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(getCoins)).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, info *dto.InfoResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrUserNotFound, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			info, err := repo.GetInfo(context.Background(), tt.userId)
			tt.expectedResp(t, info, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: sqlx.NewDb(db, "sqlmock"), cfg: DBConfig{DefaultCoins: 100}}

	tests := []struct {
		name         string
		username     string
		password     string
		mockExpect   func()
		expectedResp func(*testing.T, int, error)
	}{
		{
			name:     "success CreateUser",
			username: "testuser",
			password: "testpassword",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(insertToUsers)).
					WithArgs("testuser", "testpassword", 100).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedResp: func(t *testing.T, id int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, id)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			id, err := repo.CreateUser(context.Background(), tt.username, tt.password)
			tt.expectedResp(t, id, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
