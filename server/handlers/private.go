package handlers

import (
	"net/http"
	"order-system/handlers/api"
	"order-system/handlers/api/vendors"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

func PrivateEndpoints(e *echo.Group) {
	e.GET("/me", api.CurrentUser)
	e.POST("/orders", api.CreateOrders)
	e.PUT("/orders/:id", api.CancelOrder)
	e.GET("/orders", api.GetAllOrders)
	e.GET("/carts", api.GetCartItems)
	e.POST("/carts", api.AddItemToCart)
	e.GET("/products", api.GetAvailableProducts)
	e.POST("/orders/:id/cancel", api.CancelOrder)

	initVendorsEnpoint(e)
}

func initVendorsEnpoint(e *echo.Group) {
	vendorGroup := e.Group("/vendors", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := utils.GetCurrentUser(c)
			if user.Role != models.Vendor {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Code:    dto.ErrorUnauthorizedAccess,
					Message: "Unauthorized access",
				})
			}
			return next(c)
		}
	})

	vendorGroup.POST("/products", vendors.CreateProduct)
	vendorGroup.PUT("/products/:id", vendors.UpdateProduct)
	vendorGroup.GET("/products", vendors.GetAllVendorProducts)
	vendorGroup.GET("/products/:id/prices", vendors.GetProductPrices)
	vendorGroup.POST("/products/:id/prices", vendors.SetProductPrice)
	vendorGroup.GET("/products/:id/stocks", vendors.GetProductStocks)
	vendorGroup.POST("/products/:id/stocks", vendors.UpdateProductStock)
	vendorGroup.PUT("/orders", vendors.UpdateOrderStatus)
	vendorGroup.GET("/orders", vendors.GetAllVendorOrders)
	vendorGroup.POST("/orders/:id/cancel", api.CancelOrder)
}
