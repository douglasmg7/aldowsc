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
