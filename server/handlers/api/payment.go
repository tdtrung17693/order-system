package api

import (
	"net/http"
	"order-system/common"
	"order-system/database/payments"

	"github.com/labstack/echo/v4"
)

// GetSupportedPaymentMethods godoc
// @Summary      Get all supported payment methods
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Success      200  "Success" {object} []dto.PaymentMethod
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/payments [get]
func GetSupportedPaymentMethods(c echo.Context) error {
	res, err := payments.FindAllPaymentMethod()

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}
