package orders

import (
	"fmt"
	"order-system/common"
	"order-system/database"
	"order-system/handlers/dto"
	"order-system/models"
	"time"

	"gorm.io/gorm"
)

type VendorProduct struct {
	VendorId uint `gorm:"column:vendor_id"`
}

func CreateOrders(userId uint, orders []models.Order) error {
	db := database.GetDBInstance()

	return db.Transaction(func(tx *gorm.DB) error {
		orderItems := make(map[int]([]models.OrderItem))

		// Find another way to
		// ensure integrity between order, vendor, order item and product
		for i, order := range orders {
			vendorIds := []VendorProduct{}
			productIds := []uint{}
			for _, item := range order.Items {
				productIds = append(productIds, item.ProductID)
			}

			// ensure that all the products in the current order
			// are coming from the same vendor
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
			orderItems[i] = orders[i].Items
			// order items will be manually created
			// in order to ensure data integrity (product_price_id)
			orders[i].Items = nil
		}

		if err := tx.Create(&orders).Error; err != nil {
			return err
		}

		userCart := models.Cart{}
		if err := tx.Where("user_id = ?", userId).First(&userCart).Error; err != nil {
			return err
		}

		// order items will be created from the corresponding cart items
		orderItemsCreateQuery := `
			insert into order_items (created_at, updated_at, product_id, quantity, order_id, product_price_id)
				(select now(), now(), ci.product_id, ci.quantity, ?,  ci.product_price_id
				from cart_items ci
				inner join carts c on ci.cart_id = c.id
				where ci.product_id in (?) and c.user_id = ?)
		`

		// keep track of cart items to be deleted
		// after creating the order items
		cartItemProductIdsToClean := []uint{}

		for key, items := range orderItems {
			productIds := []uint{}
			orderId := orders[key].ID

			for _, item := range items {
				cartItemProductIdsToClean = append(cartItemProductIdsToClean, item.ProductID)
				productIds = append(productIds, item.ProductID)
			}
			fmt.Println(productIds)
			res := tx.Debug().Exec(orderItemsCreateQuery, orderId, productIds, userId)

			if res.Error != nil {
				return res.Error
			}
		}

		// clean up cart items related to created order items
		if err := tx.Where("product_id in (?) and cart_id = ?", cartItemProductIdsToClean, userCart.ID).Delete(&models.CartItem{}).
			Error; err != nil {
			return err
		}

		// for the sake of simplicity
		// all orders are in status of "PAID"
		// after creating
		orderTransactions := []models.OrderTransaction{}
		for _, order := range orders {
			newOrderTransaction := models.OrderTransaction{
				OrderID:        order.ID,
				PreviousStatus: models.OrderZeroStatus,
				Status:         models.OrderPlaced,
			}
			newOrderTransaction.CreatedAt = time.Now()
			orderTransactions = append(orderTransactions, newOrderTransaction)
			newOrderTransaction = models.OrderTransaction{
				OrderID:        order.ID,
				PreviousStatus: models.OrderPlaced,
				Status:         models.OrderPaid,
			}
			newOrderTransaction.CreatedAt = time.Now()
			orderTransactions = append(orderTransactions, newOrderTransaction)
		}

		if err := tx.Create(orderTransactions).Error; err != nil {
			return err
		}

		orderIds := []uint{}
		for _, order := range orders {
			orderIds = append(orderIds, order.ID)
		}

		// export the corresponding products
		// the product transaction entries will
		// be created based on their corresponding order items
		productTransactionsCreateQuery := `
			insert into product_transactions (created_at, type, product_id, quantity, description)
				(select now(), ?, oi.product_id, -oi.quantity, concat(concat('order ', oi.order_id), ' placed')
				from order_items oi
				inner join orders o on oi.order_id = o.id
				where o.id in (?))`

		if err := tx.Exec(productTransactionsCreateQuery, models.TransactionTypeOut, orderIds).Error; err != nil {
			return err
		}

		return nil
	})
}

func FindOrder(id uint) (dto.OrderDto, error) {
	db := database.GetDBInstance()
	order := dto.OrderDto{}

	orderQuery := `
		select o.*, pm.name as payment_method_name, u.name as vendor_name, u1.name as user_name, sum(coalesce(pp1.price,0) * oi.quantity) as total_price, ot1.created_at as status_change_time, ot1.status as status
			from orders o
			left join users u on o.vendor_id = u.id
			left join users u1 on o.user_id = u1.id
			left join order_items oi on o.id = oi.order_id
			left join product_prices pp1 on (oi.product_price_id = pp1.id)
			left join product_prices pp2 on  (oi.product_price_id = pp2.id and
											(pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))

			left join order_transactions ot1 on (o.id = ot1.order_id)
			left join order_transactions ot2 on (o.id = ot2.order_id and
												(ot1.created_at < ot2.created_at or
												(ot1.created_at = ot2.created_at and ot1.id < ot2.id)))
			left join payment_methods pm on o.payment_method_id = pm.id
			where ot2.id is null and  pp2.id is null and o.id = ?
			group by o.id, pp1.price, ot1.status, o.created_at, ot1.created_at, u.name, u1.name, pm.name
	`

	if err := db.Raw(orderQuery, id).Scan(&order).Error; err != nil {
		return order, err
	}

	orderItemsDB := []models.OrderItem{}
	if err := db.Preload("Product").Preload("ProductPrice").Where("order_id = ?", id).Find(&orderItemsDB).Error; err != nil {
		return order, err
	}

	order.Items = []dto.OrderItemDto{}
	for _, item := range orderItemsDB {
		order.Items = append(order.Items, dto.OrderItemDto{
			ProductID:   uint32(item.ProductID),
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   item.ProductPrice.Price,
		})
	}

	return order, nil
}

func OrderExists(userId uint, orderId uint, isVendor bool) (bool, error) {
	db := database.GetDBInstance()
	var exists bool
	idName := "user_id"
	if isVendor {
		idName = "vendor_id"
	}
	err := db.Model(models.Order{}).
		Select("count(*) > 0").
		Where(fmt.Sprintf("%s = ? and id = ?", idName), userId, orderId).
		Find(&exists).
		Error
	return exists, err
}

func getOrderNextStatus(status models.OrderStatus) (models.OrderStatus, error) {
	switch status {
	case models.OrderPaid:
		return models.OrderShipping, nil
	case models.OrderPlaced:
		return models.OrderPaid, nil
	case models.OrderShipping:
		return models.OrderShipped, nil
	default:
		return models.OrderZeroStatus, common.ErrorOrderFinalStateReached
	}
}

// Process the order into its next state
func SetNextStatusForOrder(orderId uint) error {
	db := database.GetDBInstance()
	status, err := FindOrderStatus(orderId)

	if err != nil {
		return err
	}

	nextStatus, err := getOrderNextStatus(status)

	if err != nil {
		return err
	}

	return db.Create(&models.OrderTransaction{
		PreviousStatus: status,
		Status:         nextStatus,
		OrderID:        orderId,
	}).Error
}

// Find all the orders that are made by an user
func FindAllOrdersOfUser(userId uint, paginationQuery dto.PaginationQuery) (*dto.PaginationResponse[dto.OrderDto], error) {
	db := database.GetDBInstance()
	o := []dto.OrderDto{}
	pageIndex := paginationQuery.PageIndex
	itemsPerPage := paginationQuery.ItemsPerPage
	filterStatus, ok := paginationQuery.Filters["status"]

	params := []interface{}{userId}
	countParams := []interface{}{userId}

	whereQuery := "where ot2.id is null and  pp2.id is null and o.user_id = ?"
	limitQuery := ""

	countQuery := `select count(o.id)
		from orders o
		left join order_transactions ot1 on (o.id = ot1.order_id)
        left join order_transactions ot2 on (o.id = ot2.order_id and
                                              (ot1.created_at < ot2.created_at or
                                               (ot1.created_at = ot2.created_at and ot1.id < ot2.id)))
		where ot2.id is null and user_id = ?`
	if ok {
		whereQuery += " and ot1.status = ?"
		countQuery += " and ot1.status = ?"
		params = append(params, filterStatus)
		countParams = append(countParams, filterStatus)
	}

	if itemsPerPage > 0 {
		params = append(params, pageIndex*itemsPerPage, itemsPerPage)
		limitQuery = " offset ? limit ?"
	}

	res := db.Raw(`
	select o.*, u.name as vendor_name, sum(coalesce(pp1.price,0) * oi.quantity) as total_price, ot1.created_at as status_change_time, ot1.status as status
		from orders o
		left join users u on o.vendor_id = u.id
        left join order_items oi on o.id = oi.order_id
        left join product_prices pp1 on (oi.product_price_id = pp1.id)
        left join product_prices pp2 on  (oi.product_price_id = pp2.id and
                                         (pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))

        left join order_transactions ot1 on (o.id = ot1.order_id)
        left join order_transactions ot2 on (o.id = ot2.order_id and
                                              (ot1.created_at < ot2.created_at or
                                               (ot1.created_at = ot2.created_at and ot1.id < ot2.id))) `+whereQuery+
		`
		 group by o.id, pp1.price, ot1.status, o.created_at, ot1.created_at, u.name 
		 order by ot1.created_at desc`+limitQuery, params...).Scan(&o)

	if res.Error != nil {
		return nil, res.Error
	}

	total := 0

	if err := db.Raw(countQuery, countParams...).Scan(&total).Error; err != nil {
		return nil, err
	}

	return &dto.PaginationResponse[dto.OrderDto]{
		Items:        o,
		Total:        total,
		PageIndex:    pageIndex,
		ItemsPerPage: itemsPerPage,
	}, nil
}

// Find all the orders that are managed by a vendor user
func FindAllOrdersOfVendor(vendorId uint, paginationQuery dto.PaginationQuery) (*dto.PaginationResponse[dto.OrderDto], error) {
	db := database.GetDBInstance()
	o := []dto.OrderDto{}
	pageIndex := paginationQuery.PageIndex
	itemsPerPage := paginationQuery.ItemsPerPage
	filterStatus, ok := paginationQuery.Filters["status"]

	params := []interface{}{vendorId}
	countParams := []interface{}{vendorId}

	whereQuery := "where ot2.id is null and  pp2.id is null and o.vendor_id = ?"
	limitQuery := ""

	countQuery := `select count(o.id)
		from orders o
		left join order_transactions ot1 on (o.id = ot1.order_id)
        left join order_transactions ot2 on (o.id = ot2.order_id and
                                              (ot1.created_at < ot2.created_at or
                                               (ot1.created_at = ot2.created_at and ot1.id < ot2.id)))
		where ot2.id is null and vendor_id = ?`
	if ok {
		whereQuery += " and ot1.status = ?"
		countQuery += " and ot1.status = ?"
		params = append(params, filterStatus)
		countParams = append(countParams, filterStatus)
	}

	if itemsPerPage > 0 {
		params = append(params, pageIndex*itemsPerPage, itemsPerPage)
		limitQuery = " offset ? limit ?"
	}

	res := db.Debug().Raw(`
	select o.*, sum(coalesce(pp1.price,0) * oi.quantity) as total_price, ot1.created_at as status_change_time, ot1.status as status
		from orders o
		left join users u on o.vendor_id = u.id
        left join order_items oi on o.id = oi.order_id
        left join product_prices pp1 on (oi.product_price_id = pp1.id)
        left join product_prices pp2 on  (oi.product_price_id = pp2.id and
                                         (pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))

        join order_transactions ot1 on (o.id = ot1.order_id)
        left join order_transactions ot2 on (o.id = ot2.order_id and
                                              (ot1.created_at < ot2.created_at or
                                               (ot1.created_at = ot2.created_at and ot1.id < ot2.id))) `+whereQuery+
		`
		 group by o.id, pp1.price, ot1.status, o.created_at, ot1.created_at
		 order by ot1.created_at desc`+limitQuery, params...).Scan(&o)

	if res.Error != nil {
		return nil, res.Error
	}

	total := 0

	if err := db.Raw(countQuery, countParams...).Scan(&total).Error; err != nil {
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
				Description: fmt.Sprintf("cancel order %d", id),
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
