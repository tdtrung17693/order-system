package utils

import (
	"order-system/handlers/dto"
	"order-system/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var user *models.User

func GetCurrentUser(c echo.Context) *models.User {
	if user == nil {
		jwtUser := c.Get("user").(*jwt.Token)
		claims := jwtUser.Claims.(*dto.JwtCustomClaims)
		user = &claims.User
	}
	return user
}
