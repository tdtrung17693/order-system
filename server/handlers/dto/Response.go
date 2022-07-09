package dto

import "errors"

type ResponseError error

var (
	ErrorEmailExist             ResponseError = errors.New("email_exists")
	ErrorPasswordMismatched     ResponseError = errors.New("password_mismatched")
	ErrorGeneric                ResponseError = errors.New("generic_error")
	ErrorUnauthorizedAccess     ResponseError = errors.New("unauthorized_access")
	ErrorInvalidCredentials     ResponseError = errors.New("invalid_credentials")
	ErrorInternalServerError    ResponseError = errors.New("internal_server_error")
	ErrorInsufficientPermission ResponseError = errors.New("insufficient_permission")
	ErrorInsufficientQuantity   ResponseError = errors.New("insufficient_stock_quantity")
)

type ErrorResponse struct {
	Code    ResponseError `json:"code"`
	Message string        `json:"message"`
}

type JSONResponse struct {
	Data interface{} `json:"data"`
}

type UserLogInResponse struct {
	AccessToken string `json:"accessToken"`
}
