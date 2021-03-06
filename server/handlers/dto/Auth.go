package dto

import (
	"order-system/models"
)

type UserLoginDto struct {
	Email    string `json:"email" valid:"required~email_required,email~email_invalid"`
	Password string `json:"password" valid:"required~password_required"`
}

type UserRegisterDto struct {
	Name            string          `json:"name" valid:"required~name_required"`
	Email           string          `json:"email" valid:"required~email_required,email~email_invalid"`
	Password        string          `json:"password" valid:"required~password_required"`
	ConfirmPassword string          `json:"confirmPassword" valid:"required~password_required"`
	Role            models.UserRole `json:"role" valid:"required~role_required"`
}

type UserDto struct {
	ID     uint            `json:"id"`
	Name   string          `json:"name"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	CartID uint            `json:"cart_id"`
}

type UserLogInResponse struct {
	AccessToken string `json:"accessToken"`
}
