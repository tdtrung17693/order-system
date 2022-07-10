package utils

import (
	"fmt"
	"order-system/config"
	"order-system/handlers/dto"
	"order-system/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user models.User) (string, error) {
	var err error

	claims := dto.JwtCustomClaims{
		User: dto.UserDto{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Role:   user.Role,
			CartID: user.Cart.ID,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	fmt.Printf("+%v", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetConfig().JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
