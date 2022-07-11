package utils

import (
	"net/http"
	"order-system/common"

	"github.com/labstack/echo/v4"
)

func BindAndValidate(c echo.Context, payload interface{}) error {
	if err := c.Bind(payload); err != nil {
		return common.ErrorInternalServerError
	}

	if err := c.Validate(payload); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}

func InternalServerError() error {
	return &echo.HTTPError{}
}
