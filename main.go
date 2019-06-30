package main

import (
	"bytes"
	"database/sql"
	"time"

	// "database/sql"
	// "encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/douglasmg7/money"
	// "github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/html/charset"
)

var db *sql.DB

// Configuration file.
type configuration struct {
	User           string      `json:"user"`
	Password       string      `json:"password"`
	FilterMinPrice money.Money `json:"filterMinPrice"`
	FilterMaxPrice money.Money `json:"filterMaxPrice"`
}

// Development mode.
var devMode bool

// Configuration.
var config configuration
var categExc []string

func init() {
	// Log file.
	_ = os.Mkdir("./log", os.ModePerm)
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

func main() {
	var err error

	// Start data base.
	db, err = sql.Open("sqlite3", "./db/aldo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

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

	// xmlFile, err := os.Open("./xml/test.xml")
	// if err != nil {
	// log.Fatalln(err)
	// }
	// defer xmlFile.Close()

	aldoXMLDoc := xmlDoc{}
	decoder := xml.NewDecoder(os.Stdin)
	// decoder := xml.NewDecoder(xmlFile)
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

	categExc = readList("list/categExc.list")
	timer := time.Now()
	err = aldoXMLDoc.process()
	fmt.Println("Time to run (s):", time.Since(timer).Seconds())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Finish.\n\n")
}

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
