package api

import (
	"errors"
	"fmt"
	"net/http"
	"order-system/database/carts"
	"order-system/database/products"
	"order-system/handlers/dto"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

func AddItemToCart(c echo.Context) error {
	o := new(dto.AddCartItemDto)

	if err := c.Bind(o); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err := c.Validate(o); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	currentUser := utils.GetCurrentUser(c)
	fmt.Println(o)
	price, err := products.FindProductLatestPrice(o.ProductID)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	if err = carts.AddItemToCart(currentUser.Cart.ID, o.ProductID, uint(o.Quantity), price.ID); err != nil {
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
