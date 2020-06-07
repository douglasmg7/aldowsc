package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

// Check consistency between aldo db and zunka db.
func checkConsistency() {
	log.Printf("Checking consistency...")
	// Get zunka products.
	zunkaProducts, err := getAllAldoZunkasiteProducts()
	if err != nil {
		log.Printf("[warn] Consistency not checked, could not get zunkasite products.")
		return
	}

	// Get db products.
	dbProducts, err := getAllDbProducts()
	if err != nil {
		log.Printf("[warn] Consistency not checked, could not get Aldo db products.")
		return
	}

	// Create db products map.
	mapDbProducts := make(map[string]aldoutil.Product)
	for _, dbProduct := range dbProducts {
		mapDbProducts[dbProduct.Code] = dbProduct
	}

	// Create zunkaiste products map.
	mapZunkaProducts := make(map[string]zunkaSiteProduct)
	for _, zunkaProduct := range zunkaProducts {
		mapZunkaProducts[zunkaProduct.MongodbId] = zunkaProduct
	}

	// Check all zunka products.
	for _, zunkaProduct := range zunkaProducts {
		// Must have a valid code, because fave dealerNmae="Aldo".
		if zunkaProduct.Code == "" {
			log.Printf("[warn] zunkasite aldo product have dealer code = \"\", zunka product _id: %s", zunkaProduct.MongodbId)
			continue
		}
		dbProduct, ok := mapDbProducts[zunkaProduct.Code]
		// Aldo product not exist for zunkasite product.
		if !ok {
			// Disable product, if active.
			if zunkaProduct.DealerProductActive {
				log.Printf("[warn] aldo product not exist for zunkasite product. zunka product _id: %s, code: %s", zunkaProduct.MongodbId, zunkaProduct.Code)
				disableZunkasiteProduct(zunkaProduct.MongodbId)
			}
			continue
		}
		if zunkaProduct.DealerProductPrice != dbProduct.DealerPrice.ToFloat64() || zunkaProduct.DealerProductActive != dbProduct.Availability {
			// Update db product with mongodbId, because zunkasite pointing to aldo product.
			if dbProduct.MongodbId == "" {
				_, err = db.Exec(stmProductUpdateMongodbIdByCode, zunkaProduct.MongodbId, dbProduct.Code)
				if !checkSQLError(err, stmProductUpdateMongodbIdByCode) {
					dbProduct.MongodbId = zunkaProduct.MongodbId
					log.Printf("[debug] Product Aldo updated. MongodbId was included. Code: %v, MongodbId: %v", dbProduct.Code, dbProduct.MongodbId)
				}
			}
			log.Printf("[warn] aldo product different from zunkasite product. zunka product _id: %s, code: %s, zunka price: %v, db price: %v, zunka active: %v, db Availability: %v", zunkaProduct.MongodbId, zunkaProduct.Code, zunkaProduct.DealerProductPrice, dbProduct.DealerPrice.ToFloat64(), zunkaProduct.DealerProductActive, dbProduct.Availability)
			updateZunkasiteProduct(&dbProduct)
		}
	}

	// Check all Aldo db products.
	for _, dbProduct := range dbProducts {
		// Check only products created at zunkasite.
		if dbProduct.MongodbId != "" {
			_, ok := mapZunkaProducts[dbProduct.MongodbId]
			if !ok {
				log.Printf("[warn] zunkasite product not exist for aldo product. Aldo product code: %s, zunka product _id: %s", dbProduct.Code, dbProduct.MongodbId)
				dbProduct.MongodbId = ""
				_, err = db.Exec(stmProductUpdateMongodbIdByCode, "", dbProduct.Code)
				if !checkSQLError(err, stmProductUpdateMongodbIdByCode) {
					log.Printf("[debug] Product Aldo updated. MongodbId was removed. Product code: %v", dbProduct.Code)
				}
			}
		}
	}
	log.Printf("Check consistency finished")
}

// Get all db products.
func getAllDbProducts() (products []aldoutil.Product, err error) {
	err = db.Select(&products, stmProductSelectAll)
	if checkSQLError(err, stmProductSelectAll) {
		return products, err
	}
	return products, nil
}

// Update zunkasite product price and availability.
func updateZunkasiteProduct(product *aldoutil.Product) error {

	// log.Printf("Runing updateZunkasiteProduct")
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
		err = errors.New(fmt.Sprintf("Error updating product on zunkaite.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		checkError(err)
		return err
	}
	log.Printf("Product updated, id: %s, active: %v, price: %v", data.ID, data.Active, data.Price)
	return nil
}

// Disable zunkasite product.
func disableZunkasiteProduct(productId string) error {
	// Product not created at zunkasite.
	if productId == "" {
		return nil
	}
	// log.Printf("product.DealerPrice.ToString(): %s", product.DealerPrice.ToString())

	// JSON data.
	data := struct {
		ID string `json:"storeProductId"`
	}{
		productId,
	}
	reqBody, err := json.Marshal(data)
	// log.Printf("reqBody: %s", reqBody)
	if checkError(err) {
		return err
	}

	// Request product update.
	client := &http.Client{}
	req, err := http.NewRequest("POST", zunkaSiteHost()+"/setup/product/disable", bytes.NewBuffer(reqBody))
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
		err = errors.New(fmt.Sprintf("Error disabling product on zunkaite.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		checkError(err)
		return err
	}
	log.Printf("Product disabled, id: %s", data.ID)
	return nil
}

// Get all aldo zunkasite products.
func getAllAldoZunkasiteProducts() (storeProducts []zunkaSiteProduct, err error) {
	// Request product update.
	client := &http.Client{}
	req, err := http.NewRequest("GET", zunkaSiteHost()+"/setup/products/aldo", nil)
	req.Header.Set("Content-Type", "application/json")
	if checkError(err) {
		return storeProducts, err
	}
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if checkError(err) {
		return storeProducts, err
	}
	defer res.Body.Close()

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		return storeProducts, err
	}
	// No 200 status.
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Error getting aldo products from zunkasite.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		checkError(err)
		return storeProducts, err
	}
	// Unmarshal.
	err = json.Unmarshal(resBody, &storeProducts)
	if checkError(err) {
		return storeProducts, err
	}

	return storeProducts, nil
}
