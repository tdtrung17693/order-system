package orders

import (
	"fmt"
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
		// Get vendor id from items

		orderItems := make(map[int]([]models.OrderItem))

		// Find another way to
		// ensure integrity between order, vendor, order item and product
		for i, order := range orders {
			vendorIds := []VendorProduct{}
			productIds := []uint{}
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
			orderItems[i] = orders[i].Items
			// order items will be manually created
			// in order to ensure data integrity (product_price_id)
			orders[i].Items = nil
		}

		if err := tx.Create(orders).Error; err != nil {
			return err
		}

		userCart := models.Cart{}
		if err := tx.Where("user_id = ?", userId).First(&userCart).Error; err != nil {
			return err
		}

		orderItemsCreateQuery := `
			insert into order_items (created_at, updated_at, product_id, quantity, order_id, product_price_id)
				(select now(), now(), ci.product_id, ci.quantity, ?,  ci.product_price_id
				from cart_items ci
				inner join carts c on ci.cart_id = c.id
				where ci.product_id in (?) and c.user_id = ?)
		`
		cartItemProductIdsToClean := []uint{}
		for key, items := range orderItems {
			productIds := []uint{}
			orderId := orders[key].ID

			for _, item := range items {
				cartItemProductIdsToClean = append(cartItemProductIdsToClean, item.ProductID)
				productIds = append(productIds, item.ProductID)
			}

			fmt.Println(productIds, orderId, userId)

			res := tx.Exec(orderItemsCreateQuery, orderId, productIds, userId)

			if res.Error != nil {
				return res.Error
			}
		}

		// // clean up cart items related to created order items
		if err := tx.Where("product_id in (?) and cart_id = ?", cartItemProductIdsToClean, userCart.ID).Delete(&models.CartItem{}).
			Error; err != nil {
			return err
		}
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

		productTransactionsCreateQuery := `
			insert into product_transactions (created_at, type, product_id, quantity, description)
				(select now(), ?, oi.product_id, -oi.quantity, concat(concat('order ', oi.order_id), ' placed')
				from order_items oi
				inner join orders o on oi.order_id = o.id
				where o.id in (?))
		`

		if err := tx.Create(orderTransactions).Error; err != nil {
			return err
		}

		orderIds := []uint{}
		for _, order := range orders {
			orderIds = append(orderIds, order.ID)
		}
		fmt.Println(orderIds)
		if err := tx.Raw(productTransactionsCreateQuery, models.TransactionTypeOut, orderIds).Error; err != nil {
			return err
		}

		return nil
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

func getOrderNextStatus(status models.OrderStatus) (models.OrderStatus, error) {
	switch status {
	case models.OrderPaid:
		return models.OrderShipping, nil
	case models.OrderPlaced:
		return models.OrderPaid, nil
	case models.OrderShipping:
		return models.OrderShipped, nil
	default:
		return models.OrderZeroStatus, dto.ErrorOrderFinalStateReached
	}
}

func SetNextStatusForOrder(orderId uint) error {
	db := database.GetDBInstance()
	status, err := FindOrderStatus(orderId)

	if err != nil {
		return err
	}

	nextStatus, err := getOrderNextStatus(status)
	fmt.Println(nextStatus)
	if err != nil {
		return err
	}

	return db.Create(&models.OrderTransaction{
		PreviousStatus: status,
		Status:         nextStatus,
		OrderID:        orderId,
	}).Error
}

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

	res := db.Raw(`
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
