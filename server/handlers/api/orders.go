package api

import (
	"bytes"
	"encoding/csv"
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
	payload := new(dto.OrdersCreateDto)

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

	currentUser := utils.GetCurrentUser(c)

	newOrders := []models.Order{}
	items := []models.OrderItem{}

	for i, order := range payload.Orders {
		newOrders = append(newOrders, models.Order{
			UserID:          currentUser.ID,
			PaymentMethodID: payload.PaymentMethodId,
			ShippingAddress: payload.RecipientAddress,
			RecipientName:   payload.RecipientName,
			RecipientPhone:  payload.RecipientPhone,
		})

		for _, item := range order.Items {
			newOrders[i].Items = append(items, models.OrderItem{
				ProductID: uint(item.ProductID),
				Quantity:  item.Quantity,
			})
		}
	}

	err := orders.CreateOrders(currentUser.ID, newOrders)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.NoContent(http.StatusOK)
}

func CancelOrder(c echo.Context) error {
	payload := new(dto.OrderUpdateStatusDto)

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

	if err = c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorGeneric,
			Message: err.Error(),
		})
	}

	if err = c.Validate(payload); err != nil {
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
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func ExportCSV(c echo.Context) error {
	status := c.QueryParam("status")

	currentUser := utils.GetCurrentUser(c)
	var filteredOrders *dto.PaginationResponse[dto.OrderDto]
	var err error

	filters := make(map[string]string)

	if len(status) > 0 {
		filters["status"] = status
	}

	if currentUser.Role == models.Vendor {
		filteredOrders, err = orders.FindAllOrdersOfVendor(currentUser.ID, dto.PaginationQuery{
			Filters: filters,
		})
	} else {
		filteredOrders, err = orders.FindAllOrdersOfUser(currentUser.ID, dto.PaginationQuery{
			Filters: filters,
		})
	}

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	b := &bytes.Buffer{}
	writer := csv.NewWriter(b)
	writer.Write(
		[]string{"Id", "Created At", "Updated At", "Total Price", "Status"},
	)
	for _, order := range filteredOrders.Items {
		writer.Write([]string{
			fmt.Sprintf("%d", order.Id),
			order.CreatedAt.Format("Jan 02 2006 15:04 -0700"),
			order.UpdatedAt.Format("Jan 02 2006 15:04 -0700"),
			order.TotalPrice.String(),
			string(order.Status),
		})
	}

	writer.Flush()
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment;filename=orders.csv")
	_, err = c.Response().Write(b.Bytes())

	return err
}
