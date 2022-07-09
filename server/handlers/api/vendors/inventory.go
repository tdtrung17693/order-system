package vendors

import (
	"net/http"
	"order-system/database/products"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetProductStocks(c echo.Context) error {
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	currentUser := utils.GetCurrentUser(c)
	product, err := products.FindProductById(uint(pId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	if product.VendorID != currentUser.ID {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Code:    dto.ErrorInsufficientPermission,
			Message: "insufficient_permission",
		})
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProductStock(c echo.Context) error {
	o := new(dto.UpdateProductStockDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

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

	product, err := products.FindProductById(uint(pId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	if product.VendorID != currentUser.ID {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrorInsufficientPermission,
			Message: "insufficient_permission",
		})
	}

	var quantity int
	if o.Type == models.TransactionTypeIn {
		err = products.ImportProductStock(uint(pId), int(o.Quantity), o.Description)
	} else {
		quantity, err = products.FindProductStockQuantity(uint(pId))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrorInternalServerError,
				Message: dto.ErrorInternalServerError.Error(),
			})
		}

		if quantity < o.Quantity {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrorInternalServerError,
				Message: dto.ErrorInternalServerError.Error(),
			})
		}

		err = products.ExportProductStock(uint(pId), int(o.Quantity), o.Description)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
