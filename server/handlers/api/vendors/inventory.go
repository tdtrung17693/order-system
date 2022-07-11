package vendors

import (
	"net/http"
	"order-system/common"
	"order-system/database/products"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UpdateProductStock godoc
// @Summary      Update product stock (import/export)
// @Tags         vendor-products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.UpdateProductStockDto false "Product stock update request"
// @Param id path int true "Product id"
// @Success      200  "Success"
// @Failure      400  "Invalid request / Insufficient stock quantity" {object} echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/products/:id/stocks [post]
func UpdateProductStock(c echo.Context) error {
	payload := new(dto.UpdateProductStockDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		return common.ErrorInternalServerError
	}

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	product, err := products.FindProductById(uint(pId))

	if err != nil {
		return common.ErrorInternalServerError
	}

	if product.VendorID != currentUser.ID {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: common.ErrorInsufficientPermission.Error(),
		}
	}

	var quantity int
	if payload.Type == models.TransactionTypeIn {
		err = products.ImportProductStock(uint(pId), int(payload.Quantity), payload.Description)
	} else {
		quantity, err = products.FindProductStockQuantity(uint(pId))

		if err != nil {
			return common.ErrorInternalServerError
		}

		if quantity < payload.Quantity {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: common.ErrorInsufficientQuantity.Error(),
			}
		}

		err = products.ExportProductStock(uint(pId), int(payload.Quantity), payload.Description)
	}

	if err != nil {
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
