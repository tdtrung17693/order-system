package models

import "time"

type TransactionType string

const (
	TransactionTypeIn  TransactionType = "in"
	TransactionTypeOut TransactionType = "out"
)

type ProductTransaction struct {
	ID          uint            `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time       `json:"createdAt"`
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	ProductID   uint            `json:"productId"`
	Quantity    int             `json:"quantity"`
}
