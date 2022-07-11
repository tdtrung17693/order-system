package vendors

import (
	"net/http"
	"order-system/common"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/services/products"
	"order-system/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CreateProduct godoc
// @Summary     Create a new product
// @Tags         vendor-products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.CreateProductDto true "Product to be created"
// @Success      200  "Success"
// @Failure      400  "Invalid request" {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/products [post]
func CreateProduct(c echo.Context) error {
	payload := new(dto.CreateProductDto)

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
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
		return common.ErrorInternalServerError
	}

	newProduct.VendorID = currentUser.ID
	return c.JSON(http.StatusOK, newProduct)
}

// UpdateProduct godoc
// @Summary     Update a product
// @Tags         vendor-products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.UpdateProductDto true "Update product request"
// @Param id path int true "Product id"
// @Success      200  "Success"
// @Failure      400  "Invalid request" {object}  echo.HTTPError
// @Failure      403  "Insufficient permission (when try to update a product that belongs other vendor)" {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/products/:id [put]
func UpdateProduct(c echo.Context) error {
	payload := new(dto.UpdateProductDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	product, err := products.FindProductById(uint(pId))

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if product.VendorID != currentUser.ID {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: common.ErrorInsufficientPermission.Error(),
		}
	}

	return products.UpdateProduct(uint(pId), *payload)
}

// GetAllVendorProducts godoc
// @Summary     Get all products that belongs to the logged in vendor
// @Tags         vendor-products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload query dto.PaginationQuery false "Pagination request"
// @Success      200  "Success" {object} dto.PaginationResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/products [get]
func GetAllVendorProducts(c echo.Context) error {
	currentUser := utils.GetCurrentUser(c)
	p := dto.ParsePaginationRequest(c)

	paginatedRes, err := products.FindProductsOfVendor(currentUser.ID, *p)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, paginatedRes)
}

// SetProductPrice godoc
// @Summary     Set the unit price of a product
// @Tags         vendor-products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload body dto.SetProductPriceDto true "Set product price request"
// @Param id path int true "Product id"
// @Success      200  "Success"
// @Failure      400  "Invalid request" {object}  echo.HTTPError
// @Failure      403  "Insufficient permission (when try to set price of a product that belongs other vendor)" {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/products/:id/prices [post]
func SetProductPrice(c echo.Context) error {
	payload := new(dto.SetProductPriceDto)
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if err := utils.BindAndValidate(c, payload); err != nil {
		return err
	}

	currentUser := utils.GetCurrentUser(c)

	product, err := products.FindProductById(uint(pId))

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	if product.VendorID != currentUser.ID {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: common.ErrorInsufficientPermission.Error(),
		}
	}

	if err := products.SetProductPrice(uint(pId), payload.Price); err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

// GetProductPrices godoc
// @Summary      Get unit price history of a product
// @Tags         vendor-products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param id path int true "Product id"
// @Success      200  "Success"
// @Failure      403  "Insufficient permission (when try to get price history of a product that belongs other vendor)" {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/vendors/products/:id/prices [get]
func GetProductPrices(c echo.Context) error {
	pIdParam := c.Param("id")

	pId, err := strconv.ParseUint(pIdParam, 10, 64)

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	currentUser := utils.GetCurrentUser(c)

	product, err := products.FindProductById(uint(pId))

	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	if product.VendorID != currentUser.ID {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: common.ErrorInsufficientPermission.Error(),
		}
	}

	prices, err := products.GetProductPrices(uint(pId))
	if err != nil {
		c.Logger().Error(err.Error())
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, prices)
}
