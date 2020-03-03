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

func Test_CreateInsertProductQuery(t *testing.T) {
	result := createStmInsertProduct(&aldoutil.Product{}, "")
	log.Println(result)
	if result == "" {
		t.Errorf("Received a empty string, want some string")
	}
}
