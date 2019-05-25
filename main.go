package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

// Aldo products.
type aldoProducts struct {
	Produto []aldoProduct `xml:"produto"`
}

// Aldo product.
type aldoProduct struct {
	Codigo    string `xml:"codigo,attr"`
	Marca     string `xml:"marca,attr"`
	Categoria string `xml:"categoria,attr"`
	Descricao string `xml:"descricao,attr"`
	// Unidade           string `xml:"unidade,attr"`
	// Multiplo          string `xml:"multiplo,attr"`		*
	Preco string `xml:"preco,attr"`
	// Precoeup          string `xml:"precoeup,attr"`
	Peso             string `xml:"peso,attr"`
	DescricaoTecnica string `xml:"descricao_tecnica,attr"`
	Disponivel       string `xml:"disponivel,attr"`
	// Ipi               string `xml:"ipi,attr"`
	Dimensoes string `xml:"dimensoes,attr"`
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
	Foto string `xml:"foto,attr"`
	// DescricaoAmigavel string `xml:"descricao_amigavel,attr"`
	// CategoriaTi       string `xml:"categoria_ti,attr"`
	TempoGarantia    string `xml:"tempo_garantia,attr"`
	ProcedimentosRma string `xml:"procedimentos_rma,attr"`
	// EmpFilial         string `xml:"emp_filial,attr"`
	// Potencia          string `xml:"potencia,attr"`
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

	products := aldoProducts{}
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel
	// decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(&products)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("products: ", products)
	// fmt.Println("product-1: ", products.Produto[1])
	fmt.Println("product-1-code: ", products.Produto[1].Codigo)
	fmt.Println("product-1-empFilial: ", products.Produto[1].EmpFilial)
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
