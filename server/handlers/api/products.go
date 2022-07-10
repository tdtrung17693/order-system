package api

import (
	"net/http"
	"order-system/database/products"
	"order-system/handlers/dto"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

func GetAvailableProducts(c echo.Context) error {
	p := dto.ParsePaginationRequest(c)

	currentUser := utils.GetCurrentUser(c)

	paginatedRes, err := products.FindAvailableProducts(currentUser.ID, *p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.JSON(http.StatusOK, paginatedRes)
}
