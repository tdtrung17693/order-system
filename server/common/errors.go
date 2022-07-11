package common

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrorEmailExist             error = errors.New("email_exists")
	ErrorPasswordMismatched     error = errors.New("password_mismatched")
	ErrorGeneric                error = errors.New("generic_error")
	ErrorUnauthorizedAccess     error = errors.New("unauthorized_access")
	ErrorInvalidCredentials     error = errors.New("invalid_credentials")
	ErrorInsufficientPermission error = errors.New("insufficient_permission")
	ErrorInsufficientQuantity   error = errors.New("insufficient_stock_quantity")
	ErrorOrderFinalStateReached error = errors.New("order_final_status_reached")
	ErrorResourceNotFound       error = errors.New("resource_not_found")
)

var (
	ErrorInternalServerError error = &echo.HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "internal_server_error",
	}
)
