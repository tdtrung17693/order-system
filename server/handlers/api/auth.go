package api

import (
	"net/http"

	"order-system/database/users"
	"order-system/handlers/dto"
	"order-system/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// type AuthController struct{}

func Login(c echo.Context) error {
	o := new(dto.UserLoginDto)

	if err := c.Bind(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	if err := c.Validate(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	storedUser, err := users.FindUserByEmail(o.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: dto.ErrorInvalidCredentials, Message: "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(o.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: dto.ErrorInvalidCredentials, Message: "Invalid credentials"})
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
	o := new(dto.UserRegisterDto)

	if err := c.Bind(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	if err := c.Validate(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	userExists, err := users.UserExists(o.Email)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	if userExists {
		return c.JSON(http.StatusOK, dto.ErrorResponse{Code: dto.ErrorEmailExist, Message: "Email exists."})
	}

	if o.ConfirmPassword != o.Password {
		return c.JSON(http.StatusOK, dto.ErrorResponse{Code: dto.ErrorPasswordMismatched, Message: "Email exists."})
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(o.Password), 10)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	// To make things simple
	// user doesn't need to veriy the email
	user, err := users.CreateUser(
		o.Name,
		o.Email,
		string(encryptedPassword),
		o.Role,
	)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: dto.ErrorInternalServerError, Message: "Internal server error."})
	}

	token, err := utils.CreateToken(user)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.UserLogInResponse{
		AccessToken: token,
	})
}

func CurrentUser(c echo.Context) error {
	user := utils.GetCurrentUser(c)
	return c.JSON(http.StatusOK, user)
}
