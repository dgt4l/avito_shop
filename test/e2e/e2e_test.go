package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/dgt4l/avito_shop/configs/avito_shop"
	"github.com/dgt4l/avito_shop/internal/avito_shop/auth"
	"github.com/dgt4l/avito_shop/internal/avito_shop/controller"
	"github.com/dgt4l/avito_shop/internal/avito_shop/dto"
	handler "github.com/dgt4l/avito_shop/internal/avito_shop/handler"
	repository "github.com/dgt4l/avito_shop/internal/avito_shop/repository/pgsql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() (*httptest.Server, func()) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logrus.Fatalf("Failed to load Config: %v", err)
	}

	logrus.Info("Config loaded successfully")

	db, err := repository.NewRepository(cfg.DBConfig)
	if err != nil {
		logrus.Fatalf("Failed to init db: %v", err)
	}

	logrus.Info("Database initialized successfully")

	authService := auth.NewAuth(cfg.AuthConfig)

	logrus.Info("Auth service initialized successfully")

	shopService := controller.NewShopService(db, authService, cfg.ServiceConfig)

	logrus.Info("Shop service initialized successfully")

	shopHandler := handler.NewShopHandler(shopService, authService, cfg.AppPort)

	logrus.Info("Shop handler initialized successfully")

	handler.RegisterRoutes(shopHandler)

	logrus.Info("Routes registered successfully")

	server := httptest.NewServer(shopHandler.GetEcho())

	logrus.Info("Test server started successfully")

	cleanupFunc := func() {
		conn, err := sqlx.Open(
			"postgres", fmt.Sprintf(
				"%s://%s:%s@%s:%s/%s?sslmode=%s",
				cfg.DBConfig.DBDriver,
				cfg.DBConfig.DBUser,
				cfg.DBConfig.DBPass,
				cfg.DBConfig.DBHost,
				cfg.DBConfig.DBPort,
				cfg.DBConfig.DBName,
				cfg.DBConfig.DBSSL,
			),
		)

		if err != nil {
			logrus.Fatalf("Failed to connect to database: %v", err)
		}
		defer conn.Close()

		_, err = conn.Exec("TRUNCATE TABLE users, items, inventory, transactions RESTART IDENTITY CASCADE;")
		if err != nil {
			logrus.Fatalf("Failed to truncate tables: %v", err)
		}

		logrus.Info("Database cleanup completed")
	}

	return server, cleanupFunc
}

func registerAndLogin(t *testing.T, server *httptest.Server, username, password string) string {
	authReq := dto.AuthRequest{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(authReq)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/api/auth", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	require.NoError(t, err)

	return authResp.Token
}

func TestAuth(t *testing.T) {
	server, cleanup := setupTestServer()
	defer server.Close()
	defer cleanup()

	t.Run("Successful registration and login", func(t *testing.T) {
		token := registerAndLogin(t, server, "testuser", "testpassword")
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		authReq := dto.AuthRequest{
			Username: "testuser",
			Password: "invalidpassword",
		}
		body, err := json.Marshal(authReq)
		require.NoError(t, err)

		resp, err := http.Post(server.URL+"/api/auth", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp dto.BadRequestResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		require.NoError(t, err)
		assert.Equal(t, "invalid password", errResp.Errors)
	})
}

func TestSendCoin(t *testing.T) {
	server, cleanup := setupTestServer()
	defer server.Close()
	defer cleanup()

	user1Token := registerAndLogin(t, server, "user1", "password1")
	_ = registerAndLogin(t, server, "user2", "password2")

	t.Run("Successful coin transfer", func(t *testing.T) {
		sendReq := dto.SendCoinRequest{
			ToUser: "user2",
			Amount: 50,
		}
		body, err := json.Marshal(sendReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/sendCoin", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+user1Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Not enough coins", func(t *testing.T) {
		sendReq := dto.SendCoinRequest{
			ToUser: "user2",
			Amount: 10000,
		}
		body, err := json.Marshal(sendReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/sendCoin", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+user1Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp dto.BadRequestResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		require.NoError(t, err)
		assert.Equal(t, "not enough coins", errResp.Errors)
	})
}

func TestGetInfo(t *testing.T) {
	server, cleanup := setupTestServer()
	defer server.Close()
	defer cleanup()

	userToken := registerAndLogin(t, server, "testuser_info", "testpassword")

	t.Run("Get user info", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/info", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+userToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info dto.InfoResponse
		err = json.NewDecoder(resp.Body).Decode(&info)
		require.NoError(t, err)

		assert.Equal(t, 1000, info.Coins)
	})
}

func TestPing(t *testing.T) {
	server, cleanup := setupTestServer()
	defer server.Close()
	defer cleanup()

	resp, err := http.Get(server.URL + "/api/ping")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "pong", string(body))
}
