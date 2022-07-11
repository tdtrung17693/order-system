package products

import (
	"fmt"
	"order-system/database"
	"order-system/handlers/dto"
	"order-system/models"
	"order-system/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func CreateProduct(product models.Product) error {
	db := database.GetDBInstance()
	res := db.Preload("Vendor").Create(&product)
	return res.Error
}

func FindProductById(id uint) (dto.ProductWithPrice, error) {
	o := dto.ProductWithPrice{}
	db := database.GetDBInstance()
	res := db.Raw(`
	select p.*, sum(coalesce(pt.quantity,0)) as stock_quantity, pp1.price as price
		from products p
        left join product_transactions pt on p.id = pt.product_id
        left join product_prices pp1 on (p.id = pp1.product_id)
        left join product_prices pp2 on  (p.id = pp2.product_id and
                                         (pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))
	where pp2.id is null and  p.id = ?
	group by p.id, pp1.price`, id).Scan(&o)

	if res.Error != nil {
		return dto.ProductWithPrice{}, res.Error
	}

	return o, nil
}

func FindProductsOfVendor(vendorId uint, paginationQuery dto.PaginationQuery) (*dto.PaginationResponse[dto.Product], error) {
	o := []dto.Product{}
	db := database.GetDBInstance()
	pageIndex := paginationQuery.PageIndex
	itemsPerPage := paginationQuery.ItemsPerPage

	err := db.Raw(`select p.*, sum(coalesce(pt.quantity, 0)) as stock_quantity, pp1.id as product_price_id, pp1.price as product_price
		from products p
		left join product_transactions pt on p.id = pt.product_id
		left join product_prices pp1 on (p.id = pp1.product_id)
		left join product_prices pp2 on  (p.id = pp2.product_id and
										(pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))
		where pp2.id is null and vendor_id = ?
		group by p.id, pp1.price, pp1.id
		order by p.created_at desc offset ? limit ?`, vendorId, pageIndex*itemsPerPage, itemsPerPage).Scan(&o).Error

	if err != nil {
		return nil, err
	}

	total := 0
	err = db.Raw(`select count(p.id) as quantity
		from products p
		where vendor_id = ?`, vendorId).Scan(&total).Error

	if err != nil {
		return nil, err
	}

	return &dto.PaginationResponse[dto.Product]{
		Items:        o,
		Total:        total,
		PageIndex:    pageIndex,
		ItemsPerPage: itemsPerPage,
	}, nil
}

func FindAvailableProducts(userId uint, paginationQuery dto.PaginationQuery) (*dto.PaginationResponse[dto.Product], error) {
	o := []dto.Product{}
	db := database.GetDBInstance()
	pageIndex := paginationQuery.PageIndex
	itemsPerPage := paginationQuery.ItemsPerPage

	err := db.Raw(`select d.* from (
			select p.*, sum(coalesce(pt.quantity, 0)) as stock_quantity, pp1.id as product_price_id, pp1.price as product_price
				from products p
				left join product_transactions pt on p.id = pt.product_id
				join product_prices pp1 on (p.id = pp1.product_id)
				left join product_prices pp2 on  (p.id = pp2.product_id and
												(pp1.created_at < pp2.created_at or (pp1.created_at = pp2.created_at and pp1.id < pp2.id)))
			where pp2.id is null
			group by p.id, pp1.price, pp1.id
		) d where stock_quantity >0 and vendor_id  != ?
		offset ? limit ?`, userId, pageIndex*itemsPerPage, itemsPerPage).Scan(&o).Error

	if err != nil {
		return nil, err
	}

	total := 0
	err = db.Raw(`select count(d.id) from (
			select p.id, sum(coalesce(pt.quantity, 0)) as stock_quantity
				from products p
				left join product_transactions pt on p.id = pt.product_id
				group by p.id
		) d where stock_quantity >0
	`).Scan(&total).Error

	if err != nil {
		return nil, err
	}

	return &dto.PaginationResponse[dto.Product]{
		Items:        o,
		Total:        total,
		PageIndex:    pageIndex,
		ItemsPerPage: itemsPerPage,
	}, nil
}

// Find the stock quantity of a product
// by summing all of its product transaction quantity
func FindProductStockQuantity(productId uint) (int, error) {
	db := database.GetDBInstance()

	total := 0
	res := db.Raw(`select sum(coalesce(pt.quantity, 0)) from products p
    left join product_transactions pt on p.id = pt.product_id
    where p.id = ?`, productId).Scan(&total)

	return total, res.Error
}

func UpdateProduct(id uint, product dto.UpdateProductDto) error {
	db := database.GetDBInstance()
	return db.Model(&models.Product{}).Where("id = ?", id).Updates(&models.Product{
		Name:        product.Name,
		Description: product.Description,
	}).Error
}

func SetProductPrice(productId uint, price decimal.Decimal) error {
	db := database.GetDBInstance()

	return db.Create(&models.ProductPrice{
		ProductID: productId,
		Price:     price,
	}).Error
}

func GetProductPrices(productId uint) ([]models.ProductPrice, error) {
	db := database.GetDBInstance()
	prices := []models.ProductPrice{}

	err := db.Where("product_id = ?", productId).Find(&prices).Error

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func FindProductLatestPrice(productId uint) (models.ProductPrice, error) {
	db := database.GetDBInstance()
	price := models.ProductPrice{}
	res := db.Where("product_id = ?", productId).Order("created_at DESC").First(&price)

	return price, res.Error
}

func ImportProductStock(productId uint, quantity int, description string) error {
	db := database.GetDBInstance()

	return db.Transaction(func(tx *gorm.DB) error {

		tableName := utils.GetModelTableName(models.ProductTransaction{})

		if err := tx.Raw(fmt.Sprintf("LOCK TABLE %s IN ACCESS EXCLUSIVE MODE;", tableName)).Error; err != nil {
			return err
		}

		err := tx.Create(&models.ProductTransaction{
			ProductID:   productId,
			Quantity:    quantity,
			Type:        models.TransactionTypeIn,
			Description: description,
		}).Error

		if err != nil {
			return err
		}
		return nil
	})
}

func ExportProductStock(productId uint, quantity int, description string) error {
	db := database.GetDBInstance()

	return db.Transaction(func(tx *gorm.DB) error {

		tableName := utils.GetModelTableName(models.ProductTransaction{})

		if err := tx.Raw(fmt.Sprintf("LOCK TABLE %s IN ACCESS EXCLUSIVE MODE;", tableName)).Error; err != nil {
			return err
		}

		err := tx.Create(&models.ProductTransaction{
			ProductID:   productId,
			Quantity:    -quantity,
			Type:        models.TransactionTypeOut,
			Description: description,
		}).Error

		if err != nil {
			return err
		}
		return nil
	})
}
