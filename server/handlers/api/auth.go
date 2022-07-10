package api

import (
	"net/http"

	"order-system/database/users"
	"order-system/handlers/dto"
	"order-system/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	payload := new(dto.UserLoginDto)

	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	storedUser, err := users.FindUserByEmail(payload.Email)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: dto.ErrorInvalidCredentials.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(payload.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: dto.ErrorInvalidCredentials, Message: dto.ErrorInvalidCredentials.Error()})
	}

	token, err := utils.CreateToken(storedUser)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.UserLogInResponse{
		AccessToken: token,
	})
}

// RegisterUser godoc
// @Summary      Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param payload body dto.UserRegisterDto true "user information"
// @Success      200  {object}  dto.UserLogInResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /register [post]
func RegisterUser(c echo.Context) error {
	payload := new(dto.UserRegisterDto)

	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	userExists, err := users.UserExists(payload.Email)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	if userExists {
		return c.JSON(http.StatusOK, dto.ErrorResponse{Code: dto.ErrorEmailExist, Message: dto.ErrorEmailExist.Error()})
	}

	if payload.ConfirmPassword != payload.Password {
		return c.JSON(http.StatusOK, dto.ErrorResponse{Code: dto.ErrorPasswordMismatched, Message: dto.ErrorPasswordMismatched.Error()})
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	// To make things simple
	// user doesn't need to veriy the email
	user, err := users.CreateUser(
		payload.Name,
		payload.Email,
		string(encryptedPassword),
		payload.Role,
	)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	token, err := utils.CreateToken(user)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	return c.JSON(http.StatusOK, dto.UserLogInResponse{
		AccessToken: token,
	})
}

func CurrentUser(c echo.Context) error {
	user := utils.GetCurrentUser(c)
	return c.JSON(http.StatusOK, user)
}
