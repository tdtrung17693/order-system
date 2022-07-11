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
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Panic("database err: ", err)
	}
}

func GetDBInstance() *gorm.DB {
	return db
}

func InitTestDB() {
	var err error
	testDBDsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PASS"))
	db, err = gorm.Open(postgres.Open(testDBDsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("storage err: ", err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalln("storage err: ", err)
	}
	sqlDb.SetMaxIdleConns(3)
	AutoMigrate(db)
	loadFixtures(db)
}

func loadFixtures(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		seedPaymentMethod(tx)
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
		regularUser := models.User{
			Name:     "regular test user",
			Email:    "regular.user@email.com",
			Password: string(encryptedPassword),
		}
		regularUser.ID = 1

		vendorUser1 := models.User{
			Name:     "vendor test user 1",
			Email:    "vendor.user.1@email.com",
			Password: string(encryptedPassword),
		}
		vendorUser1.ID = 2

		vendorUser2 := models.User{
			Name:     "vendor test user 2",
			Email:    "vendor.user.2@email.com",
			Password: string(encryptedPassword),
		}
		vendorUser2.ID = 3

		tx.Create([]models.User{regularUser, vendorUser1, vendorUser2})

		if err := tx.Create(&[]models.Cart{
			{UserID: regularUser.ID},
			{UserID: vendorUser2.ID},
			{UserID: vendorUser1.ID}}).Error; err != nil {
			return err
		}

		// create products for vendor users
		i := uint(1)
		for _, user := range []models.User{vendorUser1, vendorUser2} {
			for j := 0; j < 2; j += 1 {
				newProduct := models.Product{
					Name:     fmt.Sprintf("product %s", user.Name),
					VendorID: user.ID,
				}
				newProduct.ID = i
				productPrice := models.ProductPrice{
					ProductID: i,
					Price:     decimal.NewFromFloat(100.0),
				}
				productTransaction := models.ProductTransaction{
					Type:     models.TransactionTypeIn,
					Quantity: 5,
				}
				if err := tx.Create(&newProduct).Error; err != nil {
					return err
				}
				if err := tx.Create(&productPrice).Error; err != nil {
					return err
				}
				if err := tx.Create(&productTransaction).Error; err != nil {
					return err
				}
				i += 1
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal("load fixture error", err)
	}
}

func DropTestDB() {
	err := db.Migrator().DropTable(
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
		log.Fatal("cannot drop test tables", err)
	}
}

func AutoMigrate(db *gorm.DB) {
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

func SeedDB(db *gorm.DB) {
	seedPaymentMethod(db)
	seedDefaultUsers(db)

	fmt.Println("default user added!")
	fmt.Println("default password: password")
	fmt.Println("regular user email: email@example.com")
	fmt.Println("vendor user email: email.vendor@example.com")
}

func SeedSampleData(db *gorm.DB) {
	fmt.Println("seeding sample data...")
	if err := seedUsersAndProducts(db); err != nil {
		fmt.Println("error seeding sample data")
		fmt.Println(err)
		os.Exit(1)
	}
}

func seedDefaultUsers(db *gorm.DB) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// known user
		// 25 regular user
		newUser := models.User{
			Email:    "email@example.com",
			Name:     fmt.Sprintf("%s %s", faker.FirstName(), faker.LastName()),
			Password: string(encryptedPassword),
		}

		if err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&newUser).Error; err != nil {
			return err
		}

		cart := models.Cart{
			UserID: newUser.ID,
		}

		if err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&cart).Error; err != nil {
			return err
		}

		newVendorUser := models.User{
			Email:    "email.vendor@example.com",
			Name:     faker.DomainName(),
			Password: string(encryptedPassword),
			Role:     models.Vendor,
		}

		if err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&newVendorUser).Error; err != nil {
			return err
		}

		cart = models.Cart{
			UserID: newVendorUser.ID,
		}

		if err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&cart).Error; err != nil {
			return err
		}

		if err := seedProducts(tx, newVendorUser.ID, newVendorUser.Name); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalln("failed to populate default users")
	}
	return nil
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
		// 25 vendor user
		for i := 3; i < 25; i += 1 {
			newUser := models.User{
				Email:    faker.Email(),
				Name:     faker.DomainName(),
				Password: string(encryptedPassword),
				Role:     models.Vendor,
			}

			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&newUser).Error; err != nil {
				return err
			}

			cart := models.Cart{
				UserID: newUser.ID,
			}
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&cart).Error; err != nil {
				return err
			}
			if err := seedProducts(tx, newUser.ID, newUser.Name); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func seedProducts(db *gorm.DB, vendorId uint, vendorName string) error {
	for j := 0; j < 10; j += 1 {
		newProduct := models.Product{
			VendorID:    vendorId,
			Name:        fmt.Sprintf("(%s) product %d", vendorName, j),
			Description: fmt.Sprintf("(%s) product %d", vendorName, j),
			Unit:        "unit",
		}

		if err := db.Create(&newProduct).Error; err != nil {
			return err
		}

		productPrice := models.ProductPrice{
			ProductID: newProduct.ID,
			Price:     decimal.NewFromFloat(100.0),
		}

		if err := db.Create(&productPrice).Error; err != nil {
			return err
		}

		productTransaction := models.ProductTransaction{
			ProductID:   newProduct.ID,
			Description: "import",
			Type:        models.TransactionTypeIn,
			Quantity:    50,
		}

		if err := db.Create(&productTransaction).Error; err != nil {
			return err
		}
	}
	return nil
}
