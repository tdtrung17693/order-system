package api

import (
	"net/http"

	"order-system/common"
	"order-system/handlers/dto"
	"order-system/services/users"
	"order-system/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary      User logging in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param payload body dto.UserLoginDto true "user credentials"
// @Success      200  {object}  dto.UserLogInResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /login [post]
func Login(c echo.Context) error {
	payload := new(dto.UserLoginDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	storedUser, err := users.FindUserByEmail(payload.Email)

	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: common.ErrorInvalidCredentials,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(payload.Password)); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: common.ErrorInvalidCredentials,
		}
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
// @Failure      400  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /register [post]
func RegisterUser(c echo.Context) error {
	payload := new(dto.UserRegisterDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	userExists, err := users.UserExists(payload.Email)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if userExists {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: common.ErrorEmailExist.Error(),
		}
	}

	if payload.ConfirmPassword != payload.Password {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: common.ErrorPasswordMismatched.Error(),
		}
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
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
		return common.ErrorInternalServerError
	}

	token, err := utils.CreateToken(user)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, dto.UserLogInResponse{
		AccessToken: token,
	})
}

func CurrentUser(c echo.Context) error {
	user := utils.GetCurrentUser(c)
	return c.JSON(http.StatusOK, user)
}
