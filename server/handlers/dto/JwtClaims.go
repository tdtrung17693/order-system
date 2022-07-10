package dto

import (
	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	User UserDto
	jwt.StandardClaims
}
