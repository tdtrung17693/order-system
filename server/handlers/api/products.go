package api

import (
	"net/http"
	"order-system/common"
	"order-system/handlers/dto"
	"order-system/services/products"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

// GetAvailableProducts godoc
// @Summary      Get available products (in-stock products)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param Authorization header string true "With the bearer started"
// @Param payload query dto.PaginationQuery false "Pagination request"
// @Success      200  "Success" {object} dto.PaginationResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /api/products [get]
func GetAvailableProducts(c echo.Context) error {
	p := dto.ParsePaginationRequest(c)

	currentUser := utils.GetCurrentUser(c)

	paginatedRes, err := products.FindAvailableProducts(currentUser.ID, *p)

	if err != nil {
		c.Logger().Error(err)
		return common.ErrorInternalServerError
	}

	return c.JSON(http.StatusOK, paginatedRes)
}
