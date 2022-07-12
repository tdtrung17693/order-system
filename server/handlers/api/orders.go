package api

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"order-system/common"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/services/orders"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllOrders godoc
// @Summary      Get all orders of a current logged in user
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Success      200  "Success" {object} dto.PaginationResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/orders [get]
func GetAllOrders(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	p := dto.ParsePaginationRequest(c)

	paginatedRes, err := orders.FindAllOrdersOfUser(currentUser.ID, *p)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, paginatedRes)
}

// GetOrder godoc
// @Summary      Get an order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param id path int true "Order id"
// @Success      200  "Success" {object} dto.PaginationResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/orders [get]
func GetOrder(c echo.Context) error {
	oIdStr := c.Param("id")

	oId, err := strconv.Atoi(oIdStr)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	// Order could be get by user or vendor
	currentUser := utils.GetCurrentUser(c)
	vendorOrderExists, err := orders.OrderExists(currentUser.ID, uint(oId), true)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	if !vendorOrderExists {
		userOrderExists, err := orders.OrderExists(currentUser.ID, uint(oId), false)

		if err != nil {
			c.Logger().Error(err.Error())
			return common.ErrorInternalServerError
		}

		if !userOrderExists {
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: common.ErrorResourceNotFound.Error(),
			}
		}
	}

	order, err := orders.FindOrder(uint(oId))

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, order)
}

// CreateOrders godoc
// @Summary      Create orders based on chosen cart items
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.OrdersCreateDto true "The information of the orders to be created"
// @Success      200  "Success"
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/orders [post]
func CreateOrders(c echo.Context) error {
	payload := new(dto.OrdersCreateDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	newOrders := []models.Order{}

	for i, order := range payload.Orders {
		newOrders = append(newOrders, models.Order{
			UserID:          currentUser.ID,
			PaymentMethodID: payload.PaymentMethodId,
			ShippingAddress: payload.RecipientAddress,
			RecipientName:   payload.RecipientName,
			RecipientPhone:  payload.RecipientPhone,
		})

		for _, item := range order.Items {
			newOrders[i].Items = append(newOrders[i].Items, models.OrderItem{
				ProductID: uint(item.ProductID),
				Quantity:  item.Quantity,
			})
		}
	}

	err := orders.CreateOrders(currentUser.ID, newOrders)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

// CancelOrder godoc
// @Summary      Cancel an user's order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.OrderCancelRequest true "The information of the order to be cancelled"
// @Param id  path int true "The id of the order to be cancelled"
// @Success      200  "Success"
// @Failure      400  "invalid payload" echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/orders/:id/cancel [post]
func CancelOrder(c echo.Context) error {
	payload := new(dto.OrderCancelRequest)

	oIdStr := c.Param("id")

	oId, err := strconv.Atoi(oIdStr)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)
	exists, err := orders.OrderExists(currentUser.ID, uint(oId), false)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	if !exists {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: common.ErrorResourceNotFound.Error(),
		}
	}

	err = orders.CancelOrder(uint(oId))

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

// ExportCSV godoc
// @Summary      Export user's orders to CSV
// @Tags         orders
// @Accept       json
// @Produce      text/csv
// @Param Authorization header string true "With the bearer started"
// @Param status query string false "Status of the orders to be exported"
// @Success      200  "Success"
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/orders/export-csv [get]
func ExportCSV(c echo.Context) error {
	status := c.QueryParam("status")

	currentUser := utils.GetCurrentUser(c)
	var filteredOrders *dto.PaginationResponse[dto.OrderDto]
	var err error

	filters := make(map[string]string)

	if len(status) > 0 {
		filters["status"] = status
	}

	filteredOrders, err = orders.FindAllOrdersOfUser(currentUser.ID, dto.PaginationQuery{
		Filters: filters,
	})

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	b := &bytes.Buffer{}
	writer := csv.NewWriter(b)
	writer.Write(
		[]string{"Id", "Vendor", "Created At", "Updated At", "Total (USD)", "Status"},
	)
	for _, order := range filteredOrders.Items {
		writer.Write([]string{
			fmt.Sprintf("%d", order.Id),
			order.VendorName,
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
