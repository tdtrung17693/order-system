package vendors

import (
	"errors"
	"net/http"
	"order-system/common"
	"order-system/database/orders"
	"order-system/handlers/dto"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetAllVendorOrders godoc
// @Summary      Get all orders of current logged in vendors
// @Tags         vendor-orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload query dto.PaginationQuery false "Pagination request"
// @Success      200  "Success"
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/orders [get]
func GetAllVendorOrders(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	p := dto.ParsePaginationRequest(c)

	res, err := orders.FindAllOrdersOfVendor(currentUser.ID, *p)

	if err != nil {
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

// OrderNextStatus godoc
// @Summary      Move an order to its next status
// @Tags         vendor-orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param id path int true "Order id"
// @Success      200  "Success"
// @Failure      400  "The order has reached its final state" {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/orders/:id [put]
func OrderNextStatus(c echo.Context) error {
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	currentUser := utils.GetCurrentUser(c)
	_, err = orders.FindOrderOfVendor(uint(pId), currentUser.ID)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if err := orders.SetNextStatusForOrder(uint(pId)); err != nil {
		if errors.Is(err, common.ErrorOrderFinalStateReached) {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: common.ErrorOrderFinalStateReached.Error(),
			}
		}

		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

// CancelOrder godoc
// @Summary     Cancel an order by logged in vendor
// @Tags         vendor-orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param id path int true "Order id"
// @Param payload body dto.OrderCancelRequest true "Cancel order request payload"
// @Success      200  "Success"
// @Failure      400  "Invalid request" {object}  echo.HTTPError
// @Failure      404  "Order not found (or not belongs to current logged in vendor)" {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/orders/:id/cancel [post]
func CancelOrder(c echo.Context) error {
	payload := new(dto.OrderCancelRequest)

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)
	_, err = orders.FindOrderOfVendor(uint(id), currentUser.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: common.ErrorResourceNotFound.Error(),
			}
		}

		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	if err = c.Validate(payload); err != nil {
		return err
	}

	err = orders.CancelOrder(uint(id))

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
