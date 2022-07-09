package orders

import (
	"fmt"
	"order-system/database"
	"order-system/handlers/dto"
	"order-system/models"

	"gorm.io/gorm"
)

type VendorProduct struct {
	VendorId uint `gorm:"column:vendor_id"`
}

func CreateOrders(orders []models.Order) error {
	db := database.GetDBInstance()

	return db.Transaction(func(tx *gorm.DB) error {
		// Get vendor id from items

		productIds := []uint{}

		// Find another way to
		// ensure integrity between order, vendor, order item and product
		for i, order := range orders {
			vendorIds := []VendorProduct{}
			for _, item := range order.Items {
				productIds = append(productIds, item.ProductID)
			}

			err := tx.Raw(`
				select u.id as vendor_id from users u
				inner join products p on u.id = p.vendor_id	
				where p.id in (?)
				group by u.id
			`, productIds).Scan(&vendorIds).Error

			if err != nil {
				return err
			}

			l := len(vendorIds)
			if l == 0 || l > 1 {
				return fmt.Errorf("invalid_product_list")
			}

			orders[i].VendorID = vendorIds[0].VendorId
		}

		if err := tx.Create(orders).Error; err != nil {
			return err
		}

		orderTransactions := []models.OrderTransaction{}
		productTransactions := []models.ProductTransaction{}
		for _, order := range orders {
			orderTransactions = append(orderTransactions, models.OrderTransaction{
				OrderID:        order.ID,
				PreviousStatus: models.OrderZeroStatus,
				Status:         models.OrderPlaced,
			})
			orderTransactions = append(orderTransactions, models.OrderTransaction{
				OrderID:        order.ID,
				PreviousStatus: models.OrderPlaced,
				Status:         models.OrderPaid,
			})
			orderTransactions = append(orderTransactions, models.OrderTransaction{
				OrderID:        order.ID,
				PreviousStatus: models.OrderPaid,
				Status:         models.OrderShipping,
			})

			for _, item := range order.Items {
				productTransactions = append(productTransactions, models.ProductTransaction{
					Type:      models.TransactionTypeOut,
					Quantity:  -item.Quantity,
					ProductID: item.ProductID,
				})
			}
		}

		if err := tx.Create(orderTransactions).Error; err != nil {
			return err
		}

		return tx.Create(productTransactions).Error
	})
}

func UpdateOrderStatus(id uint, status models.OrderStatus) error {
	db := database.GetDBInstance()
	res := db.Model(&models.Order{}).Where("id = ?", id).Update("status", status)

	return res.Error
}

func FindOrderOfUser(id uint, userId uint) (models.Order, error) {
	db := database.GetDBInstance()
	o := models.Order{}
	res := db.Preload("Items").Where("user_id = ? and id = ?", id, userId).Find(&o)

	return o, res.Error
}

func FindOrderOfVendor(id uint, vendorId uint) (models.Order, error) {
	db := database.GetDBInstance()
	o := models.Order{}
	res := db.Preload("Items").Where("vendor_id = ? and id = ?", id, vendorId).Find(&o)

	return o, res.Error
}

func FindAllOrdersOfUser(userId uint, paginationQuery dto.PaginationQuery) (*dto.PaginationResponse[dto.OrderDto], error) {
	db := database.GetDBInstance()
	o := []dto.OrderDto{}
	pageIndex := paginationQuery.PageIndex
	itemsPerPage := paginationQuery.ItemsPerPage
	res := db.Raw(`
	select o.*, u.name as vendor_name, sum(coalesce(pp1.price,0)) as total_price, ot1.status as status
		from orders o
		left join users u on o.vendor_id = u.id
        left join order_items oi on o.id = oi.order_id
        join product_prices pp1 on (oi.product_price_id = pp1.id)
        left join product_prices pp2 on  (oi.product_price_id = pp2.id and
                                         (pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))

        join order_transactions ot1 on (o.id = ot1.order_id)
        left join order_transactions ot2 on (o.id = ot2.order_id and
                                              (ot1.created_at < ot2.created_at or
                                               (ot1.created_at = ot2.created_at and ot1.id < ot2.id)))
    where ot2.id is null and  pp2.id is null and o.user_id = ?
	group by o.id, pp1.price, ot1.status, u.name, o.created_at
    order by o.created_at desc
	offset ? limit ?
	`, userId, pageIndex*itemsPerPage, itemsPerPage).Scan(&o)

	if res.Error != nil {
		return nil, res.Error
	}

	total := 0
	err := db.Raw(`select count(p.id) as quantity
		from orders p
		where user_id = ?`, userId).Scan(&total).Error

	if err != nil {
		return nil, err
	}

	return &dto.PaginationResponse[dto.OrderDto]{
		Items:        o,
		Total:        total,
		PageIndex:    pageIndex,
		ItemsPerPage: itemsPerPage,
	}, nil
}

func FindAllOrdersOfVendor(vendorId uint, paginationQuery dto.PaginationQuery) (*dto.PaginationResponse[dto.OrderDto], error) {
	db := database.GetDBInstance()
	o := []dto.OrderDto{}
	pageIndex := paginationQuery.PageIndex
	itemsPerPage := paginationQuery.ItemsPerPage
	res := db.Raw(`
	select o.*, sum(coalesce(pp1.price,0)) as total_price, ot1.created_at as status_change_time, ot1.status as status
		from orders o
		left join users u on o.vendor_id = u.id
        left join order_items oi on o.id = oi.order_id
        join product_prices pp1 on (oi.product_price_id = pp1.id)
        left join product_prices pp2 on  (oi.product_price_id = pp2.id and
                                         (pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))

         join order_transactions ot1 on (o.id = ot1.order_id)
         left join order_transactions ot2 on (o.id = ot2.order_id and
                                              (ot1.created_at < ot2.created_at or
                                               (ot1.created_at = ot2.created_at and ot1.id < ot2.id)))
    where ot2.id is null and  pp2.id is null and o.vendor_id = ?
	group by o.id, pp1.price, ot1.status, o.created_at, ot1.created_at
    order by ot1.created_at desc
	offset ? limit ?
	`, vendorId, pageIndex*itemsPerPage, itemsPerPage).Scan(&o)

	if res.Error != nil {
		return nil, res.Error
	}

	total := 0
	err := db.Raw(`select count(p.id) as quantity
		from orders p
		where vendor_id = ?`, vendorId).Scan(&total).Error

	if err != nil {
		return nil, err
	}

	return &dto.PaginationResponse[dto.OrderDto]{
		Items:        o,
		Total:        total,
		PageIndex:    pageIndex,
		ItemsPerPage: itemsPerPage,
	}, nil
}

func FindOrderStatus(orderId uint) (models.OrderStatus, error) {
	db := database.GetDBInstance()
	orderTransaction := models.OrderTransaction{}
	err := db.Where("order_id = ?", orderId).Order("created_at DESC").First(&orderTransaction).Error

	return orderTransaction.Status, err
}

func CancelOrder(id uint) error {
	db := database.GetDBInstance()

	return db.Transaction(func(tx *gorm.DB) error {
		latestOrderTransaction := models.OrderTransaction{}
		err := tx.Where("order_id = ?", id).Order("created_at DESC").First(&latestOrderTransaction).Error

		if err != nil {
			return err
		}

		if latestOrderTransaction.Status == models.OrderCancelled || latestOrderTransaction.Status == models.OrderShipped {
			return nil
		}

		order := models.Order{}
		if err := tx.Preload("Items.Product").First(&order).Error; err != nil {
			return err
		}

		transactionItems := []models.ProductTransaction{}
		for _, item := range order.Items {
			transactionItems = append(transactionItems, models.ProductTransaction{
				Type:        models.TransactionTypeIn,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Description: fmt.Sprintf("cancel order number %d", id),
			})
		}

		if err := tx.Create(&transactionItems).Error; err != nil {
			return err
		}

		orderTransaction := models.OrderTransaction{
			OrderID:        id,
			PreviousStatus: models.OrderShipping,
			Status:         models.OrderCancelled,
		}

		return tx.Create(&orderTransaction).Error
	})
}
