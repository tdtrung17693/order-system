package models

type Cart struct {
	ID     uint       `json:"id"`
	Items  []CartItem `json:"items"`
	UserID uint       `json:"userId"`
}

type CartItem struct {
	ID uint `json:"id" gorm:"primaryKey"`
	BaseWithAudit
	Product        Product      `json:"product"`
	ProductID      uint         `json:"productId"`
	Quantity       uint64       `json:"quantity"`
	ProductPrice   ProductPrice `json:"productPrice"`
	ProductPriceId uint         `json:"productPriceId"`
	CartID         uint         `json:"cartId"`
	Active         bool         `json:"active"`
}
