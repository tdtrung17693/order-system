package dto

import (
	"order-system/models"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	User models.User
	jwt.StandardClaims
}
