package vendors

import (
	"fmt"
	"net/http"
	"order-system/database/orders"
	"order-system/handlers/dto"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllVendorOrders(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	p := dto.ParsePaginationRequest(c)

	res, err := orders.FindAllOrdersOfVendor(currentUser.ID, *p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateOrderStatus(c echo.Context) error {
	o := new(dto.OrderUpdateStatusDto)

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	currentUser := utils.GetCurrentUser(c)
	_, err = orders.FindOrderOfUser(uint(id), currentUser.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	if err = c.Bind(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}
	if err = c.Validate(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	err = orders.UpdateOrderStatus(uint(id), o.Status)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func CancelOrder(c echo.Context) error {
	o := new(dto.OrderUpdateStatusDto)

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	currentUser := utils.GetCurrentUser(c)
	_, err = orders.FindOrderOfVendor(uint(id), currentUser.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	if err = c.Bind(o); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	if err = c.Validate(o); err != nil {
		return err
	}

	err = orders.CancelOrder(uint(id))

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: fmt.Sprintf("Cannot cancel order %d", id),
		})
	}

	return c.NoContent(http.StatusOK)
}
