package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/douglasmg7/money"
	"golang.org/x/net/html/charset"
)

// "github.com/rogpeppe/go-charset/charset"
// _ "github.com/rogpeppe/go-charset/data"

// import (
//     "encoding/xml"
//     "golang.org/x/net/html/charset"
// )

// decoder := xml.NewDecoder(reader)
// decoder.CharsetReader = charset.NewReaderLabel
// err = decoder.Decode(&parsed)

// Configuration file.
type configuration struct {
	User           string      `json:"user"`
	Password       string      `json:"password"`
	FilterMinPrice money.Money `json:"filterMinPrice"`
	FilterMaxPrice money.Money `json:"filterMaxPrice"`
}

type xmlDoc struct {
	Products []xmlProduct `xml:"produto"`
}

type xmlProduct struct {
	Code        string `xml:"codigo,attr"`
	Brand       string `xml:"marca,attr"`
	Category    string `xml:"categoria,attr"`
	Description string `xml:"descricao,attr"`
	// Unidade           string `xml:"unidade,attr"`
	Multiplo    string `xml:"multiplo,attr"`
	DealerPrice string `xml:"preco,attr"`
	// Suggestion price to sell.
	SuggestionPrice string `xml:"precoeup,attr"`
	Weight          string `xml:"peso,attr"`
	TecDesc         string `xml:"descricao_tecnica,attr"`
	Availability    string `xml:"disponivel,attr"`
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
	PictureLinks string `xml:"foto,attr"`
	// DescricaoAmigavel string `xml:"descricao_amigavel,attr"`
	// CategoriaTi       string `xml:"categoria_ti,attr"`
	WarrantyTime string `xml:"tempo_garantia,attr"`
	RMAProcedure string `xml:"procedimentos_rma,attr"`
	// YoutubeLink         string `xml:"link_youtube,attr"`
	// EmpFilial         string `xml:"emp_filial,attr"`
	// Potencia          string `xml:"potencia,attr"`
}

type aldoProduct struct {
	Code                string
	Brand               string
	Category            string
	Description         string
	Multiple            string
	DealerPrice         money.Money
	SuggestionPrice     money.Money
	Weight              int // Peso(gr).
	TecnicalDescription string
	Availability        bool
	Dimension           dimension
	PictureLinks        []string
	WarrantyTime        string // Days.
	RMAProcedure        string // ?

}

type dimension struct {
	length int // Comprimento(mm).
	width  int // Largura(mm).
	height int // Altura(mm).
}

// Development mode.
var devMode bool

// Configuration.
var config configuration
var categExc []string

type AldoProducts []aldoProduct

func main() {

	// // Get products from last update.
	// productsRead := new(AldoProducts)
	// err := readGob("./data/products.gob", productsRead)
	// if err != nil {
	// log.Fatalln(err)
	// }
	// fmt.Printf("type: %T\n", *productsRead)
	// fmt.Println("Prod qtd: ", len(*productsRead))
	// fmt.Println("product-1", (*productsRead)[0])
	// for _, prod := range *productsRead {
	// fmt.Println(prod.Code, "\n")
	// }

	// http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c?wsdl
	// http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c
	// http://webservice.aldo.com.br/asp.net/ferramentas/saldoproduto.ashx?u=146612&p=zunk4c&codigo=20764-8&qtde=1&emp_filial=1

	// wsdl
	// url := `http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?wsdl`

	// Product quatity.
	// url := fmt.Sprintf(`http://webservice.aldo.com.br/asp.net/ferramentas/saldoproduto.ashx?u=%s&p=%s&codigo=%s&qtde=%s&emp_filial=%s`, config.WsAldo.User, config.WsAldo.Password, "20764-8", "1", "1")

	// All products.
	// url := fmt.Sprintf(`http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=%s&p=%s`, config.WsAldo.User, config.WsAldo.Password)

	// Allnations
	// url := fmt.Sprintf(`http://wspub.allnations.com.br/wsIntEstoqueClientesV2/ServicoReservasPedidosExt.asmx/RetornarListaProdutosEstoque?CodigoCliente=%s&Senha=%s&Data=%s`, config.AllNations.User, config.AllNations.Password, config.AllNations.LastReqTimeIni)

	// url := "http://www.google.com/robots.txt"

	// log.Println("url: ", url)

	// res, err := http.Get(url)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer res.Body.Close()
	// // bodyResult, err := ioutil.ReadAll(res.Header)
	// bodyResult, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for k, v := range res.Header {
	// 	log.Printf("%s - %s", k, v)
	// }
	// log.Printf("body: %s", bodyResult)

	// xmlFile, err := os.Open("./xml/arquivo_integracao_exemplo.xml")
	xmlFile, err := os.Open("./xml/test.xml")
	if err != nil {
		log.Fatalln(err)
	}

	aldoXMLDoc := xmlDoc{}
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel
	// decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(&aldoXMLDoc)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("products: ", products)
	// fmt.Println("product-1: ", products.Produto[1])
	// fmt.Println("Code: ", aldoXMLDoc.Products[1].Code)
	// fmt.Println("Description: ", aldoXMLDoc.Products[1].Description)
	// fmt.Println("Price: ", aldoXMLDoc.Products[1].Price)
	defer xmlFile.Close()

	// var f interface{}
	// decoder := xml.NewDecoder(xmlFile)
	// decoder.CharsetReader = charset.NewReader
	// err = decoder.Decode(&f)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("f: ", f)
	// defer xmlFile.Close()

	// aldoProds := aldoProducts{}
	// decoder := xml.NewDecoder(xmlFile)
	// decoder.CharsetReader = charset.NewReaderLabel
	// decoder.Decode(&aldoProds)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// sbFile, _ := ioutil.ReadAll(xmlFile)
	// log.Printf("%s", sbFile)

	// var f interface{}
	// err = xml.Unmarshal(sbFile, &f)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(f)

	// log.Println("aldo products: ", aldoProds)

	categExc = readList("list/categExc.list")
	aldoProducts := aldoXMLDoc.process()
	fmt.Println("Qutatity before serialize: ", len(aldoProducts))

	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err = e.Encode(aldoProducts)
	// err = writeGob("./data/products.gob", aldoProducts)
	if err != nil {
		log.Fatal(err)
	}

	var decodedProducts AldoProducts

	d := gob.NewDecoder(b)
	// d.CharsetReader = charset.NewReaderLabel
	for {
		err = d.Decode(&decodedProducts)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Qutatity after serialize: ", len(decodedProducts))
	// fmt.Println(decodedProducts)
	fmt.Println(decodedProducts[200])

	// categExc = readList("list/categExc.list")
	// aldoProducts := aldoXMLDoc.process()
	// fmt.Println("Qutatity before serialize: ", len(aldoProducts))

	// b := new(bytes.Buffer)
	// e := xml.NewEncoder(b)
	// err = e.Encode(aldoProducts)
	// // err = writeGob("./data/products.gob", aldoProducts)
	// if err != nil {
	// log.Fatal(err)
	// }

	// var decodedProducts AldoProducts

	// d := xml.NewDecoder(b)
	// d.CharsetReader = charset.NewReaderLabel
	// for {
	// err = d.Decode(&decodedProducts)
	// if err == io.EOF {
	// break
	// }
	// if err != nil {
	// log.Fatal(err)
	// }
	// }

	// fmt.Println("Qutatity after serialize: ", len(decodedProducts))
	// // fmt.Println(decodedProducts)
	// fmt.Println(decodedProducts[200])

	log.Printf("Finish.\n\n")
}

func init() {
	// Log file.
	logFile, err := os.OpenFile("./log/ws-aldo.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// Log configuration.
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Ldate | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// production or development mode
	setMode()

	// Configuration file.
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// sbFile, _ := ioutil.ReadAll(file)
	// log.Printf("%s", sbFile)
	config = configuration{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("WsAldo: ", config.WsAldo)
	// fmt.Println("User: ", config.WsAldo.User)
	// fmt.Println("Password: ", config.WsAldo.Password)
}

/**************************************************************************************************
* Statistics.
**************************************************************************************************/

func (doc *xmlDoc) process() (products AldoProducts) {
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
	mCategoryAllQtd := map[string]int{}
	mCategoryUsedQtd := map[string]int{}
	// var brand map[string]int
	// List of categories to get.
	// var available int
	for _, xmlProduct := range doc.Products {
		// List all categories.
		elem, _ := mCategoryAllQtd[xmlProduct.Category]
		mCategoryAllQtd[xmlProduct.Category] = elem + 1
		// Filter by categories.
		var product aldoProduct
		if !isCategorieHabToBeUsed(xmlProduct.Category) {
			prodcutQtyCutByCategFilter++
			continue
		}
		// Categories.
		product.Category = xmlProduct.Category
		// List used categories.
		elem, _ = mCategoryUsedQtd[product.Category]
		mCategoryUsedQtd[product.Category] = elem + 1
		//Price.
		var err error
		product.DealerPrice, err = money.Parse(xmlProduct.DealerPrice, ",")
		if err != nil {
			log.Printf("Could not convert price, product code: %s, price: %s\n", xmlProduct.Code, xmlProduct.DealerPrice)
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
		// Description.
		product.Description = xmlProduct.Description
		// Brands.
		product.Brand = xmlProduct.Brand

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
		product.WarrantyTime = "4 dias"
		product.RMAProcedure = "no-procedure"

		products = append(products, product)
		// Pric sum.
		priceSum += product.DealerPrice

		// fmt.Printf("[%s] - %s - R$%.2f\n", product.Category, product.Description, product.DealerPrice)
		// log.Println(product.DealerPrice)
		// log.Println()
	}
	// Average price.
	// averagePrice := priceSum.Divide(len(products))

	// log.Printf("Min price: %.2f\n", minPrice)
	// log.Printf("Max price: %.2f\n", maxPrice)
	// log.Printf("Max price code product: %s\n", maxPriceCodeProduct)
	// log.Printf("Max price desc product: %s\n", maxPriceDescriptionProduct)
	// log.Printf("Sum price: %f", priceSum)
	// log.Printf("Average price: %.4f", averagePrice)
	log.Printf("Products quantity: %d", len(doc.Products))
	log.Printf("Products quantity cut by min price(%.2f): %d", config.FilterMinPrice, prodcutQtyCutByMinPrice)
	log.Printf("Products quantity cut by max price(%.2f): %d", config.FilterMaxPrice, prodcutQtyCutByMaxPrice)
	log.Printf("Products quantity cut by categories filter: %d", prodcutQtyCutByCategFilter)
	log.Printf("Product used quantity: %d", len(products))
	log.Printf("All  Categories quantity: %d", len(mCategoryAllQtd))
	log.Printf("Used Categories quantity: %d", len(mCategoryUsedQtd))
	writeList(&mCategoryUsedQtd, "list/categUse.list")
	writeList(&mCategoryAllQtd, "list/categAll.list")
	return products
}

/**************************************************************************************************
* Encode / decode.
**************************************************************************************************/
// // writeGob encode to a binary file.
// func writeGob(filePath string, data AldoProducts) error {
// file, err := os.Create(filePath)
// if err == nil {
// encoder := xml.NewEncoder(file)
// encoder.Encode(data)
// }
// file.Close()
// return err
// }

// // readGob decode from binary file.
// func readGob(filePath string, data *AldoProducts) error {
// file, err := os.Open(filePath)
// if err == nil {
// decoder := xml.NewDecoder(file)
// err = decoder.Decode(data)
// }
// file.Close()
// return err
// }

/**************************************************************************************************
* Util.
**************************************************************************************************/
// readlist uppercase, remove spaces and create a list of lines.
func readList(fileName string) []string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Replace(string(b), " ", "", -1)
	s = strings.ToUpper(s)
	return strings.Split(s, "\n")
}

// // readlist uppercase, remove spaces and create a list of lines.
// func getCategMapToUse(fileName string) *map[string]int {
// 	b, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// s := strings.Replace(string(b), " ", "", -1)
// 	// s = strings.ToUpper(s)
// 	ss := strings.Split(string(b), "\n")
// 	m := map[string]int{}
// 	for _, s := range ss {
// 		m[s] = 1
// 	}
// 	return &m
// }

// isCategorieHabToBeUsed verify if categorie is to be used.
func isCategorieHabToBeUsed(categorie string) bool {
	categorie = strings.ToUpper(strings.Replace(categorie, " ", "", -1))
	for _, categorieExc := range categExc {
		if strings.HasPrefix(categorie, categorieExc) {
			// fmt.Printf("Prefix : %s\n", lExc)
			// fmt.Printf("Exclude: %s\n\n", l)
			return false
		}
	}
	return true
}

// writeList write a list to a file.
func writeList(m *map[string]int, fileName string) {
	b := bytes.Buffer{}
	ss := []string{}
	// Sort.
	for k, v := range *m {
		ss = append(ss, fmt.Sprintf("%s (%d)\n", k, v))
	}
	sort.Strings(ss)
	// To buffer.
	for _, s := range ss {
		b.WriteString(s)
	}
	// Write to file.
	err := ioutil.WriteFile(fileName, bytes.TrimRight(b.Bytes(), "\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

/**************************************************************************************************
* Run mode.
**************************************************************************************************/

// Define production or development mode.
func setMode() {
	for _, arg := range os.Args[1:] {
		if arg == "dev" {
			devMode = true
			log.Printf("Start - development mode.\n")
			return
		}
	}
	log.Printf("Start - production mode.\n")
}
