package dto

import "order-system/models"

type UpdateProductStockDto struct {
	Quantity    int                    `json:"quantity" valid:"required~quantity_empty,numeric~quantity_not_numeric"`
	Type        models.TransactionType `json:"type" valid:"required~update_type_empty,in(in|out)~invalid_stock_change_type"`
	Description string                 `json:"description"`
}
