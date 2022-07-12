package dto

import (
	"order-system/models"
	"time"

	"github.com/shopspring/decimal"
)

type OrdersCreateDto struct {
	Orders           []OrderCreateDto
	PaymentMethodId  string `json:"paymentMethodId"`
	RecipientAddress string `json:"recipientAddress"`
	RecipientName    string `json:"recipientName"`
	RecipientPhone   string `json:"recipientPhone"`
}

type OrderCreateDto struct {
	Items []OrderItemDto `json:"items"`
}

type OrderDto struct {
	models.BaseWithAudit
	Id                uint               `json:"id" gorm:"column:id"`
	Status            models.OrderStatus `json:"status" gorm:"column:status"`
	StatusChangeTime  time.Time          `json:"statusChangeTime" gorm:"column:status_change_time"`
	TotalPrice        decimal.Decimal    `json:"totalPrice" gorm:"column:total_price"`
	PaymentMethodID   string             `json:"paymentMethodId" gorm:"column:payment_method_id"`
	PaymentMethodName string             `json:"paymentMethodName" gorm:"column:payment_method_name"`
	ShippingAddress   string             `json:"shippingAddress" gorm:"column:shipping_address"`
	RecipientName     string             `json:"recipientName" gorm:"column:recipient_name"`
	RecipientPhone    string             `json:"recipientPhone" gorm:"column:recipient_phone"`
	VendorID          uint               `json:"vendorId" gorm:"column:vendor_id"`
	VendorName        string             `json:"vendorName" gorm:"column:vendor_name"`
	UserID            uint               `json:"userId" gorm:"column:user_id"`
	UserName          string             `json:"userName" gorm:"column:user_name"`
	Items             []OrderItemDto     `json:"items" gorm:"-"`
}

type OrderCancelRequest struct {
	Status models.OrderStatus `json:"status"`
}

type OrderItemDto struct {
	ProductID   uint32          `json:"productId"`
	ProductName string          `json:"productName"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unitPrice"`
}

type ExportCSVRequest struct {
	Status models.OrderStatus `json:"status"`
}

type OrderNextStatusRequest struct {
	OrderId uint `json:"orderId"`
}
