package models

type OrderStatus string

const (
	OrderZeroStatus = "-"
	OrderPlaced     = "PLACED"
	OrderPaid       = "PAID"
	OrderShipping   = "SHIPPING"
	OrderShipped    = "SHIPPED"
	OrderCancelled  = "CANCELLED"
)

/* Assume all orders are paid
right after they were created */
type Order struct {
	Base
	Items            []OrderItem        `json:"items"`
	UserID           uint               `json:"userId"`
	VendorID         uint               `json:"vendorId"`
	PaymentMethod    PaymentMethod      `json:"paymentMethod"`
	PaymentMethodID  string             `json:"paymentMethodId"`
	ShippingAddress  string             `json:"shippingAddress"`
	RecipientName    string             `json:"recipientName"`
	RecipientPhone   string             `json:"recipientPhone"`
	OrderTransaction []OrderTransaction `json:"orderTransaction"`
}

type OrderTransaction struct {
	BaseWithPrimaryKey
	BaseWithAudit
	PreviousStatus OrderStatus `json:"previousStatus"`
	Status         OrderStatus `json:"status"`
	OrderID        uint        `json:"orderId"`
}

type OrderItem struct {
	Base
	Product        Product      `json:"product"`
	ProductID      uint         `json:"productId"`
	Quantity       int          `json:"quantity"`
	ProductPrice   ProductPrice `json:"productPrice"`
	ProductPriceId uint         `json:"productPriceId"`
	OrderId        uint         `json:"orderId"`
}
