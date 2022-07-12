package vendors

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"order-system/common"
	"order-system/handlers/dto"
	"order-system/services/orders"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
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
	oIdStr := c.Param("id")

	oId, err := strconv.ParseUint(oIdStr, 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	currentUser := utils.GetCurrentUser(c)
	exists, err := orders.OrderExists(currentUser.ID, uint(oId), true)

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

	if err := orders.SetNextStatusForOrder(uint(oId)); err != nil {
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

	oIdStr := c.Param("id")

	oId, err := strconv.Atoi(oIdStr)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)
	exists, err := orders.OrderExists(currentUser.ID, uint(oId), true)

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

	if err = c.Validate(payload); err != nil {
		return err
	}

	err = orders.CancelOrder(uint(oId))

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

// ExportCSV godoc
// @Summary      Export user's orders to CSV
// @Tags         vendor-orders
// @Accept       json
// @Produce      text/csv
// @Param Authorization header string true "With the bearer started"
// @Param status query string false "Status of the orders to be exported"
// @Success      200  "Success"
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/orders/export-csv [get]
func ExportCSV(c echo.Context) error {
	status := c.QueryParam("status")

	currentUser := utils.GetCurrentUser(c)
	var filteredOrders *dto.PaginationResponse[dto.OrderDto]
	var err error

	filters := make(map[string]string)

	if len(status) > 0 {
		filters["status"] = status
	}

	filteredOrders, err = orders.FindAllOrdersOfVendor(currentUser.ID, dto.PaginationQuery{
		Filters: filters,
	})

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	b := &bytes.Buffer{}
	writer := csv.NewWriter(b)
	writer.Write(
		[]string{"Id", "Recipient Name", "Recipient Phone", "Address", "Created At", "Updated At", "Total (USD)", "Status"},
	)
	for _, order := range filteredOrders.Items {
		writer.Write([]string{
			fmt.Sprintf("%d", order.Id),
			order.RecipientName,
			order.RecipientPhone,
			order.ShippingAddress,
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
