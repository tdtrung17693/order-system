package api

import (
	"net/http"
	"order-system/database/products"
	"order-system/handlers/dto"

	"github.com/labstack/echo/v4"
)

func GetAvailableProducts(c echo.Context) error {
	p := dto.ParsePaginationRequest(c)

	paginatedRes, err := products.FindProducts(*p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.JSON(http.StatusOK, paginatedRes)
}
