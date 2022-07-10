package database

import (
	"fmt"
	"log"
	"os"

	"order-system/config"
	"order-system/models"
	"order-system/utils"

	"github.com/bxcodec/faker/v3"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
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
	seedDB(db)
}

func GetDBInstance() *gorm.DB {
	return db
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
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

	if err != nil {
		log.Fatalln("failed to migrate database", err)
	}
}

func seedDB(db *gorm.DB) {
	seedPaymentMethod(db)

	c := config.GetConfig()
	// seed sample data in development mode
	if c.DBSeed {
		seedUsersAndProducts(db)
	}
}

func seedPaymentMethod(db *gorm.DB) {
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

func seedUsersAndProducts(db *gorm.DB) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 25 regular user
		for i := 0; i < 25; i += 1 {
			newUser := models.User{
				Email:    faker.Email(),
				Name:     faker.DomainName(),
				Password: string(encryptedPassword),
			}
			if err := tx.Create(&newUser).Error; err != nil {
				return err
			}

			if err := tx.Create(&models.Cart{
				UserID: newUser.ID,
			}).Error; err != nil {
				return err
			}
		}

		// 25 vendor user
		for i := 0; i < 25; i += 1 {
			newUser := models.User{
				Email:    faker.Email(),
				Name:     faker.DomainName(),
				Password: string(encryptedPassword),
				Role:     models.Vendor,
			}
			if err := tx.Create(&newUser).Error; err != nil {
				return err
			}
			if err := tx.Create(&models.Cart{
				UserID: newUser.ID,
			}).Error; err != nil {
				return err
			}
			for j := 0; j < 10; j += 1 {
				newProduct := models.Product{
					VendorID:    newUser.ID,
					Name:        fmt.Sprintf("(%s) product %d", newUser.Name, j),
					Description: fmt.Sprintf("(%s) product %d", newUser.Name, j),
					Unit:        "unit",
				}

				if err := tx.Create(&newProduct).Error; err != nil {
					return err
				}

				productPrice := models.ProductPrice{
					ProductID: newProduct.ID,
					Price:     decimal.NewFromFloat(100.0),
				}

				if err := tx.Create(&productPrice).Error; err != nil {
					return err
				}

				productTransaction := models.ProductTransaction{
					ProductID:   newProduct.ID,
					Description: "import",
					Type:        models.TransactionTypeIn,
					Quantity:    50,
				}

				if err := tx.Create(&productTransaction).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
