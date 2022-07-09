package carts

import (
	"errors"
	"fmt"
	"order-system/database"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"

	"gorm.io/gorm"
)

func FindUserCart(userId uint) (dto.CartDto, error) {
	db := database.GetDBInstance()
	o := models.Cart{
		UserID: userId,
	}
	res := db.Preload("Items.Product").Preload("Items.ProductPrice").FirstOrCreate(&o)

	result := dto.CartDto{}

	for _, i := range o.Items {
		item := dto.CartItemDto{
			ID:           i.ID,
			ProductID:    i.ProductID,
			ProductName:  i.Product.Name,
			ProductPrice: i.ProductPrice.Price,
			Quantity:     uint(i.Quantity),
			VendorID:     i.Product.VendorID,
		}

		result.Items = append(result.Items, item)
	}

	return result, res.Error
}

func RemoveItemFromCart(cartId uint, itemId uint) error {
	db := database.GetDBInstance()
	res := db.Where("id = ?", cartId).Association("Items").Delete(&models.CartItem{ID: itemId})

	return res
}

func AddItemToCart(cartId uint, productId uint, requiredQuantity uint, productPriceId uint) error {
	db := database.GetDBInstance()

	type result struct {
		Result bool `gorm:"column:result"`
	}
	fmt.Println(cartId, productId)
	return db.Transaction(func(tx *gorm.DB) error {
		tableName := utils.GetModelTableName(models.CartItem{})
		var err error
		if err = tx.Raw(fmt.Sprintf("LOCK TABLE %s IN ACCESS EXCLUSIVE MODE;", tableName)).Error; err != nil {
			return err
		}
		storedCartItem := models.CartItem{}
		cartItemExisted := false
		err = tx.Where("cart_id = ? and product_id = ?", cartId, productId).First(&storedCartItem).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else if err == nil {
			cartItemExisted = true
		}

		// check stock quantity and required quantity
		o := result{}
		if err := tx.Raw(`
			select sum(stock_quantity) > sum(required_quantity) as result from (
				(select sum(coalesce(pt.quantity, 0)) as stock_quantity, 0 as required_quantity from products p
					left join product_transactions pt on p.id = pt.product_id
																	where p.id = ?
				union
				select 0, coalesce(max(quantity), 0) + ? from cart_items where cart_id = ? and product_id = ?)) d
		`, productId, cartId, requiredQuantity, productId).Scan(&o).Error; err != nil {
			return err
		}

		if !o.Result {
			return dto.ErrorInsufficientQuantity
		}

		if cartItemExisted {
			return tx.Where("cart_id = ? and product_id = ?", cartId, productId).Updates(&models.CartItem{
				Quantity: storedCartItem.Quantity + uint64(requiredQuantity),
			}).Error
		}
		newCartItem := models.CartItem{
			ProductID:      productId,
			Active:         true,
			CartID:         cartId,
			ProductPriceId: productPriceId,
			Quantity:       uint64(requiredQuantity),
		}

		return tx.Create(&newCartItem).Error
	})
}
