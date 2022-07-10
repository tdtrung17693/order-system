package dto

import "errors"

var (
	ErrorEmailExist             error = errors.New("email_exists")
	ErrorPasswordMismatched     error = errors.New("password_mismatched")
	ErrorGeneric                error = errors.New("generic_error")
	ErrorUnauthorizedAccess     error = errors.New("unauthorized_access")
	ErrorInvalidCredentials     error = errors.New("invalid_credentials")
	ErrorInternalServerError    error = errors.New("internal_server_error")
	ErrorInsufficientPermission error = errors.New("insufficient_permission")
	ErrorInsufficientQuantity   error = errors.New("insufficient_stock_quantity")
	ErrorOrderFinalStateReached error = errors.New("order_final_status_reached")
)

type ErrorResponse struct {
	Code    error  `json:"code"`
	Message string `json:"message"`
}

type JSONResponse struct {
	Data interface{} `json:"data"`
}

type UserLogInResponse struct {
	AccessToken string `json:"accessToken"`
}
