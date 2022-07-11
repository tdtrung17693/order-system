package api

import (
	"errors"
	"net/http"
	"order-system/common"
	"order-system/handlers/dto"
	"order-system/handlers/websocket"
	"order-system/services/carts"
	"order-system/services/products"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

// AddItemToCart godoc
// @Summary      Add new item to user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.AddCartItemDto true "The information of the item to be added"
// @Success      200  "Success"
// @Failure      401  "Insufficient stock quantity" {object} echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/cart [post]
func AddItemToCart(c echo.Context) error {
	payload := new(dto.AddCartItemDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	price, err := products.FindProductLatestPrice(payload.ProductID)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	if err = carts.AddItemToCart(currentUser.CartID, payload.ProductID, uint(payload.Quantity), price.ID); err != nil {
		if errors.Is(err, common.ErrorInsufficientQuantity) {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	websocket.GetHub().RefreshCart(currentUser.CartID)

	return c.NoContent(http.StatusOK)
}

// DeleteCartItem godoc
// @Summary      Delete an item in user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.DeleteCartItemDto true "The information of the cart item to be deleted"
// @Success      200  "Success"
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/cart/remove-item [post]
func DeleteCartItem(c echo.Context) error {
	payload := new(dto.DeleteCartItemDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	if err := carts.RemoveItemFromCart(currentUser.CartID, payload.ProductID); err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	websocket.GetHub().RefreshCart(currentUser.CartID)
	return c.NoContent(http.StatusOK)
}

// SetCartItemQuantity godoc
// @Summary      Change the quantity of an item in user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.SetCartItemDto true "The quantity and information of the item to be changed"
// @Success      200  "Success"
// @Failure      401  "Insufficient stock quantity" {object} echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/cart [put]
func SetCartItemQuantity(c echo.Context) error {
	payload := new(dto.SetCartItemDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	if err := carts.SetCartItemQuantity(currentUser.CartID, payload.ProductID, uint(payload.Quantity)); err != nil {
		if errors.Is(err, common.ErrorInsufficientQuantity) {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	websocket.GetHub().RefreshCart(currentUser.CartID)
	return c.NoContent(http.StatusOK)
}

// GetCartItems godoc
// @Summary      Get current cart items in user's cart
// @Tags         cart
// @Param Authorization header string true "With the bearer started"
// @Produce      json
// @Success      200  "Success"
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/cart [get]
func GetCartItems(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	cart, err := carts.FindUserCart(currentUser.ID)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, cart)
}
