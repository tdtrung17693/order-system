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

func CreateProduct(c echo.Context) error {
	payload := new(dto.CreateProductDto)

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

	newProduct := models.Product{
		Name:        payload.Name,
		Description: payload.Description,
		VendorID:    currentUser.ID,
	}

	err := products.CreateProduct(newProduct)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	newProduct.VendorID = currentUser.ID
	return c.JSON(http.StatusOK, newProduct)
}

func UpdateProduct(c echo.Context) error {
	payload := new(dto.UpdateProductDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

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

	return products.UpdateProduct(uint(pId), *payload)
}

func GetAllVendorProducts(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	p := dto.ParsePaginationRequest(c)

	paginatedRes, err := products.FindProductsOfVendor(currentUser.ID, *p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

	return c.JSON(http.StatusOK, paginatedRes)
}

func SetProductPrice(c echo.Context) error {
	payload := new(dto.SetProductPriceDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

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

	if err := products.SetProductPrice(uint(pId), payload.Price); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.NoContent(http.StatusOK)
}

func GetProductPrices(c echo.Context) error {
	payload := new(dto.SetProductPriceDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: dto.ErrorInternalServerError.Error(),
		})
	}

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

	prices, err := products.GetProductPrices(uint(pId))
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: "internal_server_error",
		})
	}

	return c.JSON(http.StatusOK, prices)
}
