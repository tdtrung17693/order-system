package api

import (
	"errors"
	"net/http"
	"order-system/database/carts"
	"order-system/database/products"
	"order-system/handlers/dto"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

func AddItemToCart(c echo.Context) error {
	payload := new(dto.AddCartItemDto)

	if err := c.Bind(payload); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err := c.Validate(payload); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	currentUser := utils.GetCurrentUser(c)

	price, err := products.FindProductLatestPrice(payload.ProductID)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	if err = carts.AddItemToCart(currentUser.CartID, payload.ProductID, uint(payload.Quantity), price.ID); err != nil {
		if errors.Is(err, dto.ErrorInsufficientQuantity) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrorGeneric,
				Message: err.Error(),
			})
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.NoContent(http.StatusOK)
}

func DeleteCartItem(c echo.Context) error {
	payload := new(dto.DeleteCartItemDto)

	if err := c.Bind(payload); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err := c.Validate(payload); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	currentUser := utils.GetCurrentUser(c)

	if err := carts.RemoveItemFromCart(currentUser.CartID, payload.ProductID); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.NoContent(http.StatusOK)
}

func SetCartItemQuantity(c echo.Context) error {
	payload := new(dto.SetCartItemDto)

	if err := c.Bind(payload); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err := c.Validate(payload); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	currentUser := utils.GetCurrentUser(c)

	if err := carts.SetCartItemQuantity(currentUser.CartID, payload.ProductID, uint(payload.Quantity)); err != nil {
		if errors.Is(err, dto.ErrorInsufficientQuantity) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrorGeneric,
				Message: err.Error(),
			})
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.NoContent(http.StatusOK)

}

func GetCartItems(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	cart, err := carts.FindUserCart(currentUser.ID)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.JSON(http.StatusOK, cart)
}
