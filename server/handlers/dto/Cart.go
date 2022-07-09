package dto

import "github.com/shopspring/decimal"

type AddCartItemDto struct {
	ProductID uint   `json:"productId"`
	Quantity  uint64 `json:"quantity"`
}

type CartDto struct {
	Items []CartItemDto `json:"items"`
}

type CartItemDto struct {
	ID           uint            `json:"id"`
	ProductName  string          `json:"productName"`
	ProductID    uint            `json:"productId"`
	ProductPrice decimal.Decimal `json:"productPrice"`
	Quantity     uint            `json:"quantity"`
	VendorID     uint            `json:"vendorId"`
}
