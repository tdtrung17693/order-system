package dto

import (
	"github.com/shopspring/decimal"
)

type CreateProductDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	VendorID    uint   `json:"vendorId"`
}

type Product struct {
	ID             uint            `json:"id" gorm:"column:id"`
	Name           string          `json:"name" gorm:"column:name"`
	Description    string          `json:"description" gorm:"column:description"`
	VendorID       uint            `json:"vendorId" gorm:"column:vendor_id"`
	Unit           string          `json:"unit" gorm:"column:uint"`
	StockQuantity  uint            `json:"stockQuantity" gorm:"column:stock_quantity"`
	ProductPriceId uint            `json:"productPriceId" gorm:"column:product_price_id"`
	ProductPrice   decimal.Decimal `json:"productPrice" gorm:"column:product_price"`
}

type ProductWithPrice struct {
	Product
	Price decimal.Decimal `json:"price" gorm:"column:price"`
}

type UpdateProductDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SetProductPriceDto struct {
	Price decimal.Decimal `json:"price"`
}
