package api

import (
	"fmt"
	"net/http"
	"order-system/database/orders"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllOrders(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	p := dto.ParsePaginationRequest(c)

	paginatedRes, err := orders.FindAllOrdersOfUser(currentUser.ID, *p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, paginatedRes)
}

func CreateOrders(c echo.Context) error {
	o := new(dto.OrdersCreateDto)

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

	currentUser := utils.GetCurrentUser(c)

	newOrders := []models.Order{}
	items := []models.OrderItem{}

	for i, order := range o.Orders {
		newOrders = append(newOrders, models.Order{
			UserID:          currentUser.ID,
			VendorID:        order.VendorID,
			PaymentMethodID: o.PaymentMethodId,
			ShippingAddress: o.ShippingAddress,
			RecipientName:   o.RecipientName,
			RecipientPhone:  o.RecipientPhone,
		})

		for _, item := range order.Items {
			newOrders[i].Items = append(items, models.OrderItem{
				ProductID:      uint(item.ProductID),
				Quantity:       item.Quantity,
				ProductPriceId: item.ProductPriceID,
			})
		}
	}

	err := orders.CreateOrders(newOrders)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.NoContent(http.StatusOK)
}

func CancelOrder(c echo.Context) error {
	o := new(dto.OrderUpdateStatusDto)

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	currentUser := utils.GetCurrentUser(c)
	_, err = orders.FindOrderOfUser(uint(id), currentUser.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: dto.ErrorGeneric, Message: "Order not found"})
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
