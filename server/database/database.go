package database

import (
	"log"
	"os"

	"order-system/config"
	"order-system/models"
	"order-system/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func InitDB() {
	var err error
	dsn := config.GetDBDsn()
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("database err: ", err)
	}
	autoMigrate(db)
	initPaymentMethod(db)
}

func GetDBInstance() *gorm.DB {
	return db
}

// func TestDB() *gorm.DB {
// 	db, err := gorm.Open(, "./../realworld_test.db")
// 	if err != nil {
// 		fmt.Println("storage err: ", err)
// 	}
// 	db.DB().SetMaxIdleConns(3)
// 	db.LogMode(false)
// 	return db
// }

// func DropTestDB() error {
// 	if err := os.Remove("./../realworld_test.db"); err != nil {
// 		return err
// 	}
// 	return nil
// }

//TODO: err check
func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductTransaction{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderTransaction{},
		&models.PaymentMethod{},
		&models.ProductPrice{},
		&models.Cart{},
		&models.CartItem{},
	)
}

func initPaymentMethod(db *gorm.DB) {
	payments := []models.PaymentMethod{
		{
			ID:   "payment_cod",
			Name: "cash_on_delivery",
		},
		{
			ID:   "payment_credit",
			Name: "credit_card",
		},
	}

	err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(payments).Error

	if err != nil {
		utils.LogErrorLn("failed to populate payment methods")
		utils.LogErrorLn(err)
		os.Exit(1)
	}
}
