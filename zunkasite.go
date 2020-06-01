package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/douglasmg7/aldoutil"
)

// Zunka site product.
type zunkaSiteProduct struct {
	MongodbId           string  `json:"id"`
	Code                string  `json:"dealerProductId"`
	DealerProductActive bool    `json:"dealerProductActive"`
	DealerProductPrice  float64 `json:"dealerProductPrice"`
}

// Update zunkasite product price and availability.
func updateZunkasiteProduct(product *aldoutil.Product) error {
	// Product not created at zunkasite.
	if product.MongodbId == "" {
		return nil
	}
	// log.Printf("product.DealerPrice.ToString(): %s", product.DealerPrice.ToString())

	// JSON data.
	data := struct {
		ID     string `json:"storeProductId"`
		Active bool   `json:"dealerProductActive"`
		Price  string `json:"dealerProductPrice"`
	}{
		product.MongodbId,
		product.Availability,
		product.DealerPrice.ToString(),
	}
	reqBody, err := json.Marshal(data)
	// log.Printf("reqBody: %s", reqBody)
	if checkError(err) {
		return err
	}

	// Request product update.
	client := &http.Client{}
	req, err := http.NewRequest("POST", zunkaSiteHost()+"/setup/product/update", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if checkError(err) {
		return err
	}
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if checkError(err) {
		return err
	}
	// res, err := http.Post("http://localhost:3080/setup/product/add", "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	if checkError(err) {
		return err
	}

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		return err
	}
	// No 200 status.
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Error ao solicitar a criação do produto no servidor zunka.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		checkError(err)
		return err
	}
	return nil
}

// Get all aldo zunkasite products.
func getAllAldoZunkasiteProducts() (err error, storeProducts []zunkaSiteProduct) {
	// Request product update.
	client := &http.Client{}
	req, err := http.NewRequest("GET", zunkaSiteHost()+"/setup/products/aldo", nil)
	req.Header.Set("Content-Type", "application/json")
	if checkError(err) {
		return err, storeProducts
	}
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if checkError(err) {
		return err, storeProducts
	}
	defer res.Body.Close()

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		return err, storeProducts
	}
	// No 200 status.
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Error getting aldo products from zunkasite.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		checkError(err)
		return err, storeProducts
	}
	// Unmarshal.
	err = json.Unmarshal(resBody, &storeProducts)
	if checkError(err) {
		return err, storeProducts
	}

	return nil, storeProducts
}
