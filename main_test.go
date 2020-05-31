package main

import (
	"log"
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
}

func shutdownTest() {
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
	log.Println(result)
	if result == "" {
		t.Errorf("Received a empty string, want some string")
	}
}

// Create satment insert product.
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
