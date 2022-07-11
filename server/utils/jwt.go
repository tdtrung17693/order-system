package utils

import (
	"errors"
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

func VerifyAndParseJWT(tokenStr string) (*dto.UserDto, error) {
	// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
	claims := dto.JwtCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.GetConfig().JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims.User, nil
}
