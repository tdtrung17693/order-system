package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductPrice struct {
	ID        uint            `json:"id" gorm:"primarykey"`
	CreatedAt time.Time       `json:"createdAt"`
	DeletedAt gorm.DeletedAt  `json:"deletedAt"`
	ProductID uint            `json:"productId"`
	Product   Product         `json:"product"`
	Price     decimal.Decimal `json:"price" gorm:"type:numeric;"`
}
