package main

import (
	"database/sql"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/douglasmg7/aldoutil"
	"github.com/douglasmg7/money"
)

type xmlProduct struct {
	Code        string `xml:"codigo,attr"`
	Brand       string `xml:"marca,attr"`
	Category    string `xml:"categoria,attr"`
	Description string `xml:"descricao,attr"`
	Unit        string `xml:"unidade,attr"`
	Multiplo    string `xml:"multiplo,attr"`
	DealerPrice string `xml:"preco,attr"`
	// Suggestion price to sell.
	SuggestionPrice      string `xml:"precoeup,attr"`
	Weight               string `xml:"peso,attr"`
	TechnicalDescription string `xml:"descricao_tecnica,attr"`
	Availability         string `xml:"disponivel,attr"`
	// Ipi               string `xml:"ipi,attr"`
	Measurements string `xml:"dimensoes,attr"`
	// Abnt              string `xml:"abnt,attr"`
	// Ncm               string `xml:"ncm,attr"`
	// Origem            string `xml:"origem,attr"`
	// Ppb               string `xml:"ppb,attr"`
	// Portariappb       string `xml:"portariappb,attr"`
	// Mpdobem           string `xml:"mpdobem,attr"`
	// Dataportariappb   string `xml:"dataportariappb,attr"`
	// Icms              string `xml:"icms,attr"`
	// Reducao           string `xml:"reducao,attr"`
	// Precocomst        string `xml:"precocomst,attr"`
	// Produtocomst      string `xml:"produtocomst,attr"`
	PictureLink string `xml:"foto,attr"`
	// DescricaoAmigavel string `xml:"descricao_amigavel,attr"`
	// CategoriaTi       string `xml:"categoria_ti,attr"`
	WarrantyTime string `xml:"tempo_garantia,attr"`
	RMAProcedure string `xml:"procedimentos_rma,attr"`
	// YoutubeLink         string `xml:"link_youtube,attr"`
	// EmpFilial         string `xml:"emp_filial,attr"`
	// Potencia          string `xml:"potencia,attr"`
}

type xmlDoc struct {
	Products []xmlProduct `xml:"produto"`
}

func (doc *xmlDoc) process() (err error) {
	// Price.
	var minPrice money.Money
	minPrice = math.MaxFloat32
	var maxPrice money.Money
	// var maxPriceCodeProduct string
	// var maxPriceDescriptionProduct string
	var priceSum money.Money
	var prodcutQtyCutByMaxPrice int
	var prodcutQtyCutByMinPrice int
	var prodcutQtyCutByCategFilter int
	var productQtyCutByError int
	mCategoryAll := map[string]int{}
	mCategoryUse := map[string]int{}
	// var brand map[string]int
	// List of categories to get.
	// var available int

	var totalProductQtd int
	var usedProductQtd int

	for _, xmlProduct := range doc.Products {
		totalProductQtd++
		// List all categories.
		elem, _ := mCategoryAll[xmlProduct.Category]
		mCategoryAll[xmlProduct.Category] = elem + 1
		// Filter by categories.
		if !isCategorieSelected(xmlProduct.Category) {
			prodcutQtyCutByCategFilter++
			continue
		}
		product := aldoutil.Product{}
		// Categories.
		product.Category = xmlProduct.Category
		// List used categories.
		elem, _ = mCategoryUse[product.Category]
		mCategoryUse[product.Category] = elem + 1

		// Price.
		var err error
		product.DealerPrice, err = money.Parse(xmlProduct.DealerPrice, ",")
		if err != nil {
			log.Printf("Could not convert dealer price, product code: %s, price: %s\n", xmlProduct.Code, xmlProduct.DealerPrice)
			continue
		}

		// Suggestion price.
		product.SuggestionPrice, err = money.Parse(xmlProduct.SuggestionPrice, ",")
		if err != nil {
			log.Printf("Could not convert suggestion price, product code: %s, price: %s\n", xmlProduct.Code, xmlProduct.SuggestionPrice)
			continue
		}

		// Filter max price.
		if product.DealerPrice > config.FilterMaxPrice {
			prodcutQtyCutByMaxPrice++
			continue
		}
		// Filter min price.
		if product.DealerPrice < config.FilterMinPrice {
			prodcutQtyCutByMinPrice++
			continue
		}

		// Code.
		product.Code = xmlProduct.Code

		// Brands.
		product.Brand = xmlProduct.Brand

		// Description.
		product.Description = xmlProduct.Description

		// Unit.
		product.Unit = xmlProduct.Unit

		// Multiple (multiple of unit).
		multipleInt64, err := strconv.ParseInt(xmlProduct.Multiplo, 10, 0)
		if err != nil {
			log.Printf("Product with code %s not imported (invalid multiple), err: %s", product.Code, err)
			productQtyCutByError++
			continue
		}
		product.Multiple = int(multipleInt64)

		// Techincal description.
		product.TechnicalDescription = xmlProduct.TechnicalDescription

		// Availability.
		if strings.ToLower(strings.TrimSpace(xmlProduct.Availability)) == "sim" {
			product.Availability = true
		}
		// Weight, remove ".", change "," to "." and parse.
		weightKg, err := strconv.ParseFloat(strings.Replace(strings.ReplaceAll(xmlProduct.Weight, ".", ""), ",", ".", 1), 64)
		if err != nil {
			log.Printf("Product with code %s not imported (invalid weight), err: %s", product.Code, err)
			productQtyCutByError++
			continue
		}
		// Convert to grams.
		product.Weight = int(weightKg * 1000)

		// Get length, width and height.
		re := regexp.MustCompile(`\d*\.?\d+`)
		dims := re.FindAllString(xmlProduct.Measurements, 3)

		// Not have all dimensions.
		if len(dims) != 3 {
			log.Printf("Product with code %s not imported (invalid dimensions), err: %s", product.Code, "Not have all dimensions")
			productQtyCutByError++
			continue
		}

		// Height.
		heightCm, err := strconv.ParseFloat(dims[0], 64)
		if err != nil {
			log.Printf("Product with code %s not imported (invalid height), err: %s", product.Code, err)
			productQtyCutByError++
			continue
		}
		product.Height = int(heightCm * 10)

		// Width.
		widthCm, err := strconv.ParseFloat(dims[1], 64)
		if err != nil {
			log.Printf("Product with code %s not imported (invalid width), err: %s", product.Code, err)
			productQtyCutByError++
			continue
		}
		product.Width = int(widthCm * 10)

		// Length.
		lengthCm, err := strconv.ParseFloat(dims[2], 64)
		if err != nil {
			log.Printf("Product with code %s not imported (invalid length), err: %s", product.Code, err)
			productQtyCutByError++
			continue
		}
		product.Length = int(lengthCm * 10)

		// Picture.
		product.PictureLink = xmlProduct.PictureLink

		// Warrant.
		re = regexp.MustCompile(`\d+`)
		warrantTime64, err := strconv.ParseInt(re.FindAllString(xmlProduct.WarrantyTime, 1)[0], 10, 0)
		product.WarrantyPeriod = int(warrantTime64)
		if err != nil {
			log.Printf("Product with code %s not imported (invalid warranty time), err: %s", product.Code, err)
			productQtyCutByError++
			continue
		}

		// RMA procedure.
		product.RMAProcedure = xmlProduct.RMAProcedure

		// fmt.Println("DealerPrice: ", product.DealerPrice)
		// Max price.
		if product.DealerPrice > maxPrice {
			maxPrice = product.DealerPrice
			// maxPriceCodeProduct = product.Code
			// maxPriceDescriptionProduct = product.Description
		}

		// Min price.
		if product.DealerPrice < minPrice {
			minPrice = product.DealerPrice
		}

		// Pric sum.
		priceSum += product.DealerPrice
		// Product will be used.
		usedProductQtd++

		// Get product from db.
		dbProduct := aldoutil.Product{}
		err = dbProduct.FindByCode(db, product.Code)

		// Error.
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}

		// New product.
		if err == sql.ErrNoRows {
			// log.Println("Inserting:", product.Code)
			product.New = true
			product.CreatedAt = time.Now()
			product.ChangedAt = product.CreatedAt
			// fmt.Println("Inserted product:", product.Code)
			err = product.Save(db)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		// Product already exist.
		// fmt.Println("Product found on db:", dbProduct.Code)

		// Product changed.
		if product.Diff(&dbProduct) {
			// fmt.Println("productDb change:", dbProduct.Changed)
			// fmt.Println("productDb CreatedAt:", dbProduct.CreatedAt)
			// fmt.Println("productDb ChangedAt:", dbProduct.ChangedAt)
			// Save product history.
			err = dbProduct.SaveHistory(db)
			if err != nil {
				log.Fatal(err)
			}
			// Update changed product.
			product.Id = dbProduct.Id
			product.CreatedAt = dbProduct.CreatedAt
			product.ChangedAt = time.Now()
			product.Changed = true
			err = product.Update(db)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Product changed", product.Code)
		}
	}

	// Save all categories on db.
	// Average price.
	// averagePrice := priceSum.Divide(len(products))

	// log.Printf("Min price: %.2f\n", minPrice)
	// log.Printf("Max price: %.2f\n", maxPrice)
	// log.Printf("Max price code product: %s\n", maxPriceCodeProduct)
	// log.Printf("Max price desc product: %s\n", maxPriceDescriptionProduct)
	// log.Printf("Sum price: %f", priceSum)
	// log.Printf("Average price: %.4f", averagePrice)
	log.Printf("Products total: %d", totalProductQtd)
	log.Printf("Products cut by min price(%.2f): %d", config.FilterMinPrice, prodcutQtyCutByMinPrice)
	log.Printf("Products cut by max price(%.2f): %d", config.FilterMaxPrice, prodcutQtyCutByMaxPrice)
	log.Printf("Products cut by categories filter: %d", prodcutQtyCutByCategFilter)
	log.Printf("Products cut by error: %d", productQtyCutByError)
	log.Printf("Product used: %d", usedProductQtd)
	log.Printf("All  Categories: %d", len(mCategoryAll))
	log.Printf("Used Categories: %d", len(mCategoryUse))
	writeList(&mCategoryUse, "list/categUse.list")
	writeList(&mCategoryAll, "list/categAll.list")
	return err
}
