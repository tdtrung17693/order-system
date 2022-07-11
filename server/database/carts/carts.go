package carts

import (
	"errors"
	"fmt"
	"order-system/common"
	"order-system/database"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"

	"gorm.io/gorm"
)

func FindUserCart(userId uint) (dto.CartDto, error) {
	db := database.GetDBInstance()
	c := models.Cart{}
	err := db.Where("user_id = ?", userId).Find(&c).Error
	if err != nil {
		return dto.CartDto{}, err
	}

	o := []models.CartItem{}
	res := db.Preload("Product").Preload("ProductPrice").Preload("Product.Vendor").Where("cart_id = ?", c.ID).Order("updated_at DESC").Find(&o)

	result := dto.CartDto{
		Items: make([]dto.CartItemDto, 0),
	}

	for _, i := range o {
		item := dto.CartItemDto{
			ID:           i.ID,
			ProductID:    i.ProductID,
			ProductName:  i.Product.Name,
			ProductPrice: i.ProductPrice.Price,
			Quantity:     uint(i.Quantity),
			VendorID:     i.Product.VendorID,
			VendorName:   i.Product.Vendor.Name,
		}

		result.Items = append(result.Items, item)
	}

	return result, res.Error
}

// Simply remove an item out of the user cart
func RemoveItemFromCart(cartId uint, productId uint) error {
	db := database.GetDBInstance()
	res := db.Where("product_id = ? and cart_id = ?", productId, cartId).Delete(&models.CartItem{})

	return res.Error
}

// Add a number of product item into the cart
// The quantity will be checked and ensure that
// it does not exceed the product stock quantity
func AddItemToCart(cartId uint, productId uint, requiredQuantity uint, productPriceId uint) error {
	db := database.GetDBInstance()

	type Result struct {
		Result bool `json:"result" gorm:"column:result"`
	}

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
		o := Result{}
		if err := tx.Raw(`
			select sum(stock_quantity) > sum(required_quantity) as result from (
				(select sum(coalesce(pt.quantity, 0)) as stock_quantity, 0 as required_quantity from products p
					left join product_transactions pt on p.id = pt.product_id
																	where p.id = ?
				union
				select 0 as stock_quantity, coalesce(max(quantity), 0) + ? as required_quantity from cart_items where cart_id = ? and product_id = ?)) d
		`, productId, requiredQuantity, cartId, productId).Scan(&o).Error; err != nil {
			return err
		}

		if !o.Result {
			return common.ErrorInsufficientQuantity
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

// Set the quantity of an entry in the cart
// The quantity will be checked and ensure that
// it does not exceed the product stock quantity
func SetCartItemQuantity(cartId uint, productId uint, requiredQuantity uint) error {
	db := database.GetDBInstance()

	type Result struct {
		Result bool `json:"result" gorm:"column:result"`
	}

	return db.Transaction(func(tx *gorm.DB) error {
		tableName := utils.GetModelTableName(models.CartItem{})
		var err error
		if err = tx.Raw(fmt.Sprintf("LOCK TABLE %s IN ACCESS EXCLUSIVE MODE;", tableName)).Error; err != nil {
			return err
		}

		// check stock quantity and required quantity
		o := Result{}
		if err := tx.Raw(`
			select sum(stock_quantity) > sum(required_quantity) as result from (
				(select sum(coalesce(pt.quantity, 0)) as stock_quantity, 0 as required_quantity from products p
					left join product_transactions pt on p.id = pt.product_id
																	where p.id = ?
				union
				select 0 as stock_quantity, ? as required_quantity from cart_items where cart_id = ? and product_id = ?)) d
		`, productId, requiredQuantity, cartId, productId).Scan(&o).Error; err != nil {
			return err
		}

		if !o.Result {
			return common.ErrorInsufficientQuantity
		}

		updatedCartItem := models.CartItem{
			Quantity: uint64(requiredQuantity),
		}
		return tx.Where("cart_id = ? and product_id = ?", cartId, productId).Updates(&updatedCartItem).Error
	})
}
