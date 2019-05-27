package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	WsAldo struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"wsAldo"`
	AllNations struct {
		User           string `json:"user"`
		Password       string `json:"password"`
		LastReqTimeIni string `json:"lastReqTimeIni"`
	} `json:"allNations"`
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
	Multiplo string `xml:"multiplo,attr"`
	Price    string `xml:"preco,attr"`
	// Precoeup          string `xml:"precoeup,attr"`
	Weight       string `xml:"peso,attr"`
	TecDesc      string `xml:"descricao_tecnica,attr"`
	Availability string `xml:"disponivel,attr"`
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
	// EmpFilial         string `xml:"emp_filial,attr"`
	// Potencia          string `xml:"potencia,attr"`
}

type aldoProduct struct {
	Code                string
	Brand               string
	Category            string
	Description         string
	Multiple            string
	Price               float32
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

func main() {
	// http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c?wsdl
	// http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c
	// http://webservice.aldo.com.br/asp.net/ferramentas/saldoproduto.ashx?u=146612&p=zunk4c&codigo=20764-8&qtde=1&emp_filial=1

	// wsdl
	// url := `http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?wsdl`

	// Product quatity.
	// url := fmt.Sprintf(`http://webservice.aldo.com.br/asp.net/ferramentas/saldoproduto.ashx?u=%s&p=%s&codigo=%s&qtde=%s&emp_filial=%s`, config.WsAldo.User, config.WsAldo.Password, "20764-8", "1", "1")

	// All products.
	url := fmt.Sprintf(`http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=%s&p=%s`, config.WsAldo.User, config.WsAldo.Password)

	// Allnations
	// url := fmt.Sprintf(`http://wspub.allnations.com.br/wsIntEstoqueClientesV2/ServicoReservasPedidosExt.asmx/RetornarListaProdutosEstoque?CodigoCliente=%s&Senha=%s&Data=%s`, config.AllNations.User, config.AllNations.Password, config.AllNations.LastReqTimeIni)

	// url := "http://www.google.com/robots.txt"

	// log.Println("url: ", url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	// bodyResult, err := ioutil.ReadAll(res.Header)
	bodyResult, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range res.Header {
		log.Printf("%s - %s", k, v)
	}
	log.Printf("body: %s", bodyResult)

	xmlFile, err := os.Open("./xml/arquivo_integracao_exemplo.xml")
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
	fmt.Println("Code: ", aldoXMLDoc.Products[1].Code)
	fmt.Println("Description: ", aldoXMLDoc.Products[1].Description)
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
	products := aldoXMLDoc.process()

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

func (doc *xmlDoc) process() (products []aldoProduct) {
	// Price.
	var minPrice float32
	minPrice = math.MaxFloat32
	var maxPrice float32
	var priceSum float32
	var averagePrice float32
	var brand map[string]int
	var category map[string]int
	var available int
	for _, xmlProduct := range doc.Products {
		var product aldoProduct
		product.Brand = xmlProduct.Brand
		//Price.
		pirce, err := convertStrBrDecimalToFloat32(xmlProduct.Price)
		if err != nil {
			log.Println("Could not convert price, product code: %s, price: %s", xmlProduct.Code, xmlProduct.Price)
			continue
		}
		product.Price = pirce
		// Max price.
		if product.Price > maxPrice {
			maxPrice = product.Price
		}
		// Min price.
		if product.Price < minPrice {
			minPrice = product.Price
		}
		products = append(products, product)
		// Pric sum.
		priceSum += product.Price
	}
	averagePrice := priceSum / len(products)
	return products
}

/**************************************************************************************************
* Util.
**************************************************************************************************/
func convertStrBrDecimalToFloat32(str string) (val float32, err error) {
	str = strings.Replace(str, ".", "", -1)
	str = strings.Replace(str, ",", ".", 0)
	val64, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return float32(val64), nil
	}
	return 0, err
}

/**************************************************************************************************
* Run mode.
**************************************************************************************************/

// Define production or development mode.
func setMode() {
	for _, arg := range os.Args[1:] {
		if arg == "dev" {
			devMode = true
			log.Println("Starting - development mode.")
			return
		}
	}
	log.Println("Starting - production mode.")
}
