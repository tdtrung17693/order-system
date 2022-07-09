package utils

import (
	"order-system/config"
	"order-system/handlers/dto"
	"order-system/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user models.User) (string, error) {
	var err error

	claims := dto.JwtCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetConfig().JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
