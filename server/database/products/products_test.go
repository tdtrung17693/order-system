package products_test

import (
	"order-system/database"
	"order-system/database/products"
	"order-system/models"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func TestProductPrice(t *testing.T) {
	database.InitTestDB()
	defer database.DropTestDB()

	productId := 1
	expectedPrice := decimal.NewFromFloat(100)

	price, err := products.FindProductLatestPrice(uint(productId))
	if err != nil {
		t.Error("error while getting product price")
	}

	if !price.Price.Equal(expectedPrice) {
		t.Logf("expected: +%v", expectedPrice)
		t.Errorf("actual: +%v", price.Price)
	}
}

func TestProductStockQuantity(t *testing.T) {
	database.InitTestDB()
	defer database.DropTestDB()

	productId := 1
	expectedPrice := decimal.NewFromFloat(100)

	price, err := products.FindProductLatestPrice(uint(productId))
	if err != nil {
		t.Error("error while getting product price")
	}

	if !price.Price.Equal(expectedPrice) {
		t.Logf("expected: +%v", expectedPrice)
		t.Errorf("actual: +%v", price.Price)
	}
}

func TestImportProduct(t *testing.T) {
	database.InitTestDB()
	defer database.DropTestDB()

	productId := 1
	expectedQuantity := 10
	expectedDescription := "test 1"

	if err := products.ImportProductStock(uint(productId), 10, "test 1"); err != nil {
		t.Error("error while importing product")
	}

	testDb := database.GetDBInstance()
	productTransaction := models.ProductTransaction{}

	if err := testDb.Where("product_id = ?", uint(productId)).
		Order("created_at desc").
		First(&productTransaction).Error; err != nil {
		t.Error("error while getting product transaction", err)
	}

	if productTransaction.Description != expectedDescription ||
		productTransaction.Quantity != expectedQuantity ||
		productTransaction.Type != models.TransactionTypeIn {
		t.Log("expected: v", []interface{}{expectedQuantity, expectedDescription, models.TransactionTypeIn})
		t.Error("actual: v", []interface{}{productTransaction.Quantity, productTransaction.Description, productTransaction.Type})
	}
}

func TestExportProduct(t *testing.T) {
	database.InitTestDB()
	defer database.DropTestDB()

	productId := 1
	expectedQuantity := 10
	expectedDescription := "test 1"

	if err := products.ExportProductStock(uint(productId), 10, "test 1"); err != nil {
		t.Error("error while exporting product")
	}

	testDb := database.GetDBInstance()
	productTransaction := models.ProductTransaction{}

	if err := testDb.Where("product_id = ?", uint(productId)).
		Order("created_at desc").
		First(&productTransaction).Error; err != nil {
		t.Error("error while getting product transaction", err)
	}

	if productTransaction.Description != expectedDescription ||
		productTransaction.Quantity != -expectedQuantity ||
		productTransaction.Type != models.TransactionTypeOut {
		t.Log("expected: v", []interface{}{-expectedQuantity, expectedDescription, models.TransactionTypeOut})
		t.Error("actual: v", []interface{}{productTransaction.Quantity, productTransaction.Description, productTransaction.Type})
	}
}

func TestMain(m *testing.M) {
	cwd, _ := os.Getwd()
	godotenv.Load(path.Join(cwd, "..", "..", ".env.testing"))
	code := m.Run()
	os.Exit(code)
}
