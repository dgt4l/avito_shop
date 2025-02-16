package auth

import (
	"testing"
	"time"

	"github.com/dgt4l/avito_shop/internal/avito_shop/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestServiceAuth_GenerateToken(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	user := &models.User{
		Id:       1,
		Username: "testuser",
		Password: "testpassword",
	}

	token, err := service.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SigningKey), nil
	})
	assert.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(user.Id), claims["id"])
	assert.Equal(t, user.Username, claims["username"])
	assert.Equal(t, user.Password, claims["password"])
	assert.InDelta(t, time.Now().Add(TokenTTL).Unix(), claims["expires_at"].(float64), 1)
}

func TestServiceAuth_ParseToken_ValidToken(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	user := &models.User{
		Id:       1,
		Username: "testuser",
		Password: "testpassword",
	}

	token, err := service.GenerateToken(user)
	assert.NoError(t, err)

	userId, err := service.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.Id, userId)
}

func TestServiceAuthParseTokenInvalidKey(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         1,
		"expires_at": time.Now().Add(TokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte("wrong-key"))
	assert.NoError(t, err)

	_, err = service.ParseToken(tokenString)
	assert.Error(t, err)
}

func TestServiceAuthParseTokenInvalidToken(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	_, err := service.ParseToken("invalid-token")
	assert.Error(t, err)
}

func TestServiceAuthParseTokenExpiredToken(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         1,
		"expires_at": time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(cfg.SigningKey))
	assert.NoError(t, err)

	_, err = service.ParseToken(tokenString)
	assert.Error(t, err)
}

func TestServiceAuthParseTokenMissingClaim(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
	})
	tokenString, err := token.SignedString([]byte(cfg.SigningKey))
	assert.NoError(t, err)

	_, err = service.ParseToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, ErrClaimIdFails, err)
}

func TestServiceAuthParseTokenInvalidClaimType(t *testing.T) {
	cfg := AuthConfig{
		SigningKey: "test-key",
	}
	service := NewAuth(cfg)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       "not-a-number",
		"username": "testuser",
	})
	tokenString, err := token.SignedString([]byte(cfg.SigningKey))
	assert.NoError(t, err)

	_, err = service.ParseToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, ErrClaimIdFails, err)
}
