package auth

import (
	"errors"
	"time"

	"github.com/dgt4l/avito_shop/internal/avito_shop/models"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	GenerateToken(user *models.User) (string, error)
	ParseToken(tokenString string) (int, error)
}

type ServiceAuth struct {
	cfg AuthConfig
}

func NewAuth(cfg AuthConfig) *ServiceAuth {
	return &ServiceAuth{
		cfg: cfg,
	}
}

func (s *ServiceAuth) GenerateToken(user *models.User) (string, error) {
	logrus.Info("id:", user.Id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"id":         user.Id,
		"username":   user.Username,
		"password":   user.Password,
		"expires_at": time.Now().Add(TokenTTL).Unix(),
	})

	return token.SignedString([]byte(s.cfg.SigningKey))
}

func (s *ServiceAuth) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.cfg.SigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(float64)
		if !ok {
			return 0, errors.New("claim parsing id fails")
		}
		return int(id), nil
	}
	return 0, errors.New("claim missing")
}
