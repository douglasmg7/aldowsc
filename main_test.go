package main

import (
	"os"
	"testing"

	"github.com/douglasmg7/aldoutil"
)

func TestMain(m *testing.M) {

	setupTest()
	code := m.Run()
	shutdownTest()

	os.Exit(code)
}

func setupTest() {
	initSql3DB()
}

func shutdownTest() {
	closeSql3DB()
}

// Create satment insert product.
func Test_CreateStmInsertProduct(t *testing.T) {
	result := createStmInsert(&aldoutil.Product{}, "")
	// log.Println(result)
	if result == "" {
		t.Errorf("Received a empty string, want some string")
	}
}

// Create statement update product.
func Test_CreateStmUpdateProductByCode(t *testing.T) {
	result := createStmUpdateByCode(&aldoutil.Product{}, "")
	// log.Println(result)
	if result == "" {
		t.Errorf("Received a empty string, want some string")
	}
}

// Update zunkasite product.
func Test_DisableZunkasiteProduct(t *testing.T) {
	err := disableZunkasiteProduct("5ec52855711e8f07336c6697")
	if err != nil {
		t.Errorf("Failed. %s", err)
	}
}

// Update zunkasite product.
func Test_UpdateZunkasiteProduct(t *testing.T) {
	product := aldoutil.Product{
		MongodbId:    "5ec52855711e8f07336c6697",
		Availability: true,
		DealerPrice:  123456,
	}
	err := updateZunkasiteProduct(&product)
	if err != nil {
		t.Errorf("Failed. %s", err)
	}
}

// Get zunkasite aldo products.
func Test_GetAllAldoZunkasiteProduct(t *testing.T) {
	zunkaProducts, err := getAllAldoZunkasiteProducts()
	if err != nil {
		t.Errorf("Failed. %s", err)
	}
	// Some product.
	// log.Printf("zunkaProducts: %+v", zunkaProducts)
	if len(zunkaProducts) == 0 {
		t.Errorf("Received no zunkasite aldo products len = 0.")
	}
	// MongodbId.
	if len(zunkaProducts[0].MongodbId) != 24 {
		t.Errorf("Invalid MongodbId: %v", zunkaProducts[0].MongodbId)
	}
	// Code.
	if zunkaProducts[0].Code == "" {
		t.Errorf("Invalid code: %v", zunkaProducts[0].Code)
	}
	// Price.
	if zunkaProducts[0].DealerProductPrice < 100 {
		t.Errorf("Invalid price: %v", zunkaProducts[0].DealerProductPrice)
	}
}

// Get zunkasite aldo products.
func Test_GetAllDbProducts(t *testing.T) {
	// Get all products.
	products, err := getAllDbProducts()
	if err != nil {
		t.Errorf("Failed. %s", err)
	}
	// Some product.
	if len(products) == 0 {
		t.Errorf("Received db products len = 0")
	}
	// Code.
	if products[0].Code == "" {
		t.Errorf("Invalid code: %v", products[0].Code)
	}
	// Price.
	if products[0].DealerPrice < 100 {
		t.Errorf("Invalid price: %v", products[0].DealerPrice)
	}
}

// Check consistency.
func Test_CheckConsistency(t *testing.T) {
	checkConsistency()
}
