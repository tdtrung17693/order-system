package utils

import (
	"order-system/handlers/dto"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetCurrentUser(c echo.Context) *dto.UserDto {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(*dto.JwtCustomClaims)
	return &claims.User
}
