package handlers

import (
	"net/http"
	"order-system/handlers/api"

	"github.com/labstack/echo/v4"
)

func PublicEndpoints(e *echo.Group) {
	e.POST("/login", api.Login)
	e.POST("/register", api.RegisterUser)
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
