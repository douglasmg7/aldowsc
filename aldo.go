package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var totalChecked = 0

// Update all zunka aldo products stock..
func updateAllZunkaAldoProductsStock() {
	log.Printf("Updating zunka aldo product quantity...")

	wg := sync.WaitGroup{}

	// Get zunka products.
	allProducts, err := getAllAldoZunkasiteProducts()
	if checkError(err) {
		log.Printf("[warn] zunka aldo products quantity not updated, could not get zunkasite products.")
		return
	}

	// Products with stock.
	countHasStock := 0
	for _, product := range allProducts {
		if product.DealerProductActive && !product.SuccessfullyProcessed && product.StoreProductQtd > 0 {
			countHasStock++
		}
	}
	log.Printf("Quantity of products with stock: %d", countHasStock)

	// Check all zunka products.
	countTry := 0
	countProductsToProcess := 0
	// Try seveal times, until all products processed.
	for {
		countProductsToProcess = 0
		// Process products.
		for i, product := range allProducts {
			if product.DealerProductActive && !product.SuccessfullyProcessed {
				countProductsToProcess++
				// log.Printf("Product, Code: %s, MongodbId: %s", product.Code, product.MongodbId)
				wg.Add(1)
				go updateZunkaAldoProductStock(allProducts, i, &wg)
			}
		}
		countTry++
		if countProductsToProcess != 0 {
			log.Printf("[%d] %d products to check stock quantity", countTry, countProductsToProcess)
		}
		wg.Wait()

		// Exit if fineshed.
		if countProductsToProcess == 0 {
			goto Exit
		}
	}
Exit:
	log.Printf("Updating zunka aldo product finished")
}

// Update zunka aldo product stock.
func updateZunkaAldoProductStock(products []zunkaSiteProduct, item int, wg *sync.WaitGroup) {
	// log.Printf("checking product, code: %s, MongodbId: %s", product.Code, product.MongodbId)
	defer wg.Done()
	// Try 3 products.
	has, ok := checkAldoProductQuantity(products[item].Code, 3)
	if !ok {
		// log.Println("Not ok")
		return
	}
	if ok && has {
		if products[item].StoreProductQtd != 3 {
			updateZunkasiteProductQuantity(products[item], 3)
		}
		products[item].SuccessfullyProcessed = true
		totalChecked++
		log.Printf("Product %s checked, quantity: %d, item [%d] ", products[item].Code, 3, totalChecked)
		return
	}
	// Try 1 products.
	has, ok = checkAldoProductQuantity(products[item].Code, 1)
	if ok {
		if has {
			if products[item].StoreProductQtd != 1 {
				updateZunkasiteProductQuantity(products[item], 1)
			}
			products[item].SuccessfullyProcessed = true
			totalChecked++
			log.Printf("Product %s checked, quantity: %d, item [%d] ", products[item].Code, 1, totalChecked)
		} else {
			// Stock 0.
			if products[item].StoreProductQtd > 0 {
				updateZunkasiteProductQuantity(products[item], 0)
			}
			products[item].SuccessfullyProcessed = true
			totalChecked++
			log.Printf("Product %s checked, quantity: %d, item [%d] ", products[item].Code, 0, totalChecked)
		}
	}
}

// Update zunkasite product quantity.
func updateZunkasiteProductQuantity(product zunkaSiteProduct, newQuantity int) {
	// JSON data.
	data := struct {
		ID       string `json:"_id"`
		Quantity int    `json:"storeProductQtd"`
	}{
		product.MongodbId,
		newQuantity,
	}
	reqBody, err := json.Marshal(data)
	// log.Printf("reqBody: %s", reqBody)
	if checkError(err) {
		return
	}

	// Request product update.
	client := &http.Client{}
	req, err := http.NewRequest("POST", zunkaSiteHost()+"/setup/product/quantity", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if checkError(err) {
		return
	}
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if checkError(err) {
		return
	}
	// res, err := http.Post("http://localhost:3080/setup/product/add", "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	if checkError(err) {
		return
	}

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		return
	}
	// No 200 status.
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Error updating product %s to quantity %d, on zunkaite server.\nstatus: %v\nbody: %v", product.MongodbId, newQuantity, res.StatusCode, string(resBody)))
		checkError(err)
		return
	}
	log.Printf("Product %s was updated, _id: %s, old quantity: %d, new quantity %d", product.Code, product.MongodbId, product.StoreProductQtd, newQuantity)
	return
}

// Check aldo product quantity.
// return (has, ok).
func checkAldoProductQuantity(aldoProductID string, quantity int) (bool, bool) {
	// Not sell the last one.
	quantity = quantity + 1
	// Request.
	url := fmt.Sprintf("%s?u=%s&p=%s&codigo=%s&qtde=%d&emp_filial=%s", aldoHost, aldoUser, aldoPass, aldoProductID, quantity, "1")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	// req.Header.Set("Content-Type", "application/json")
	if checkError(err) {
		return false, false
	}
	// req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if checkError(err) {
		return false, false
	}
	defer res.Body.Close()

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		return false, false
	}
	// No 200 status.
	if res.StatusCode != 200 {
		// // err = errors.New(fmt.Sprintf("Error getting aldo products quantity.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		// err = errors.New(fmt.Sprintf("Getting aldo products quantity, product: %s, quantity: %d, status: %v", aldoProductID, quantity, res.StatusCode))
		// checkError(err)
		return false, false
	}
	// log.Printf("resBody: %s", string(resBody))
	if strings.ToLower(string(resBody)) == "sim" {
		return true, true
	}
	return false, true
}
