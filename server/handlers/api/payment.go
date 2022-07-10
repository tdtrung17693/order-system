package api

import (
	"net/http"
	"order-system/database/payments"
	"order-system/handlers/dto"

	"github.com/labstack/echo/v4"
)

func GetSupportedPaymentMethods(c echo.Context) error {
	res, err := payments.FindAllPaymentMethod()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}
