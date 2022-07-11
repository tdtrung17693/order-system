package handlers

import (
	"errors"
	"net/http"
	"order-system/handlers/api"
	"order-system/handlers/api/vendors"
	"order-system/models"
	"order-system/utils"

	"github.com/labstack/echo/v4"
)

func PrivateEndpoints(e *echo.Group) {
	e.GET("/me", api.CurrentUser)
	e.POST("/orders", api.CreateOrders)
	e.PUT("/orders/:id", api.CancelOrder)
	e.GET("/orders", api.GetAllOrders)
	e.GET("/cart", api.GetCartItems)
	e.POST("/cart", api.AddItemToCart)
	e.PUT("/cart", api.SetCartItemQuantity)
	e.POST("/cart/remove-item", api.DeleteCartItem)
	e.GET("/products", api.GetAvailableProducts)
	e.POST("/orders/:id/cancel", api.CancelOrder)
	e.GET("/orders/export-csv", api.ExportCSV)
	e.GET("/payment-methods", api.GetSupportedPaymentMethods)

	initVendorsEnpoint(e)
}

func initVendorsEnpoint(e *echo.Group) {
	vendorGroup := e.Group("/vendors", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := utils.GetCurrentUser(c)
			if user.Role != models.Vendor {
				return &echo.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: errors.New("ABC"),
				}
			}
			return next(c)
		}
	})

	vendorGroup.POST("/products", vendors.CreateProduct)
	vendorGroup.PUT("/products/:id", vendors.UpdateProduct)
	vendorGroup.GET("/products", vendors.GetAllVendorProducts)
	vendorGroup.GET("/products/:id/prices", vendors.GetProductPrices)
	vendorGroup.POST("/products/:id/prices", vendors.SetProductPrice)
	vendorGroup.POST("/products/:id/stocks", vendors.UpdateProductStock)
	vendorGroup.GET("/orders", vendors.GetAllVendorOrders)
	vendorGroup.PUT("/orders/:id", vendors.OrderNextStatus)
	vendorGroup.POST("/orders/:id/cancel", vendors.CancelOrder)
}
