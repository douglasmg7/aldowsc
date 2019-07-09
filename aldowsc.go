package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"golang.org/x/net/html/charset"
)

var db *sql.DB
var dbAldo *sqlx.DB

// File names.
var dbFileName string = "aldowsc.db"

// Paths.
var appPath, logPath, dbPath, xmlPath, listPath string

// Files.
var logFile, dbFile string

// Min and max price filter.
var maxPriceFilter, minPriceFilter int

// Development mode.
var production bool

// Categories selected to use.
var categSel []string

func init() {
	appPath := os.Getenv("ZUNKAPATH")
	logPath = path.Join(appPath, "log")
	dbPath = path.Join(appPath, "db")
	listPath = path.Join(appPath, "list")
	xmlPath = path.Join(appPath, "xml")

	os.MkdirAll(logPath, os.ModePerm)
	os.MkdirAll(dbPath, os.ModePerm)
	os.MkdirAll(listPath, os.ModePerm)
	os.MkdirAll(xmlPath, os.ModePerm)

	dbFile = path.Join(dbPath, "aldowsc.db")
	logFile = path.Join(logPath, "aldowsc.log")

	// Log file.
	// os.MkdirAll("path" + "//log/aldowsc.log")
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// logFile, err := os.OpenFile("./log/ws-aldo.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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

	// Run mode.
	if os.Getenv("ZUNKAENV") == "PRODUCTION" {
		production = true
		log.Println("Running in production mode")
	}

	// Config.
	viper.SetDefault("aldowsc.filter.minPrice", 2000)
	viper.SetDefault("aldowsc.filter.maxPrice", 100000)
	viper.SetConfigName("config")
	viper.AddConfigPath(os.Getenv("ZUNKAPATH"))
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Error reading config file: %s \n", err))
	}
	minPriceFilter = viper.GetInt("aldowsc.filter.minPrice")
	maxPriceFilter = viper.GetInt("aldowsc.filter.maxPrice")
}

func main() {
	var err error

	// Db start.
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Init db Aldo.
	dbAldo = sqlx.MustConnect("sqlite3", dbFile)

	// Read selected categories.
	log.Println("Reading selected categories list...")
	categSel = readList(path.Join(listPath, "categSel.list"))

	// Remove no more selected products.
	rmProductsNotSel()

	// Load xml file.
	log.Println("Loading and decoding xml file...")
	timer := time.Now()
	aldoXMLDoc := xmlDoc{}
	decoder := xml.NewDecoder(os.Stdin)

	// Decoding xml file.
	timer = time.Now()
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&aldoXMLDoc)
	if err != nil {
		log.Fatalln("Error decoding xml file:", err)
	}
	log.Printf("Time loading and decoding xml file: %fs", time.Since(timer).Seconds())

	// Processing products.
	log.Println("Processing products...")
	timer = time.Now()
	err = aldoXMLDoc.process()
	log.Printf("Time processing products: %fs", time.Since(timer).Seconds())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Finish.\n\n")
}

/**************************************************************************************************
* Util.
**************************************************************************************************/
// readlist lowercase, remove spaces and create a list of lines.
func readList(fileName string) []string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Replace(string(b), " ", "", -1)
	s = strings.ToLower(s)
	return strings.Split(s, "\n")
}

// writeList write a list to a file.
func writeList(m *map[string]int, fileName string) {
	b := bytes.Buffer{}
	ss := []string{}
	// Sort.
	for k, v := range *m {
		ss = append(ss, fmt.Sprintf("%s (%d)\n", strings.ToLower(k), v))
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

// Verify if categorie is to be used.
func isCategorieSelected(categorie string) bool {
	categorie = strings.ToLower(strings.Replace(categorie, " ", "", -1))
	for _, categItem := range categSel {
		if strings.HasPrefix(categorie, categItem) {
			// fmt.Printf("Prefix : %s\n", lExc)
			// fmt.Printf("Exclude: %s\n\n", l)
			return true
		}
	}
	return false
}

// Remove no more selected products from db.
func rmProductsNotSel() {
	// Get distinct categories from products on db.
	dbCategs := []string{}
	err := dbAldo.Select(&dbCategs, "SELECT distinct category FROM product")
	if err != nil {
		log.Fatal(fmt.Errorf("Get distinct categories from db: %v", err))
	}
	// Categories to be removed.
	categToRem := []string{}
	for _, dbCateg := range dbCategs {
		if !isCategorieSelected(dbCateg) {
			categToRem = append(categToRem, `"`+dbCateg+`"`)
		}
	}
	if len(categToRem) > 0 {
		log.Printf("Removing no more selected categorie(s): %s.", strings.Join(categToRem, ", "))
	}

	// Copy products to remove to history.
	tx := dbAldo.MustBegin()
	tx.MustExec(fmt.Sprintf("INSERT INTO product_history SELECT * FROM product WHERE category IN (%s)", strings.Join(categToRem, ",")))
	// Delete copied products.
	tx.MustExec(fmt.Sprintf("DELETE FROM product WHERE category IN (%s)", strings.Join(categToRem, ",")))
	err = tx.Commit()
	if err != nil {
		log.Fatal(fmt.Errorf("Removing products from db: %v", err))
	}

	// Get products to remove.
	// products := []aldoutil.Product{}
	// fmt.Printf("SELECT * FROM product WHERE category IN (%s)", strings.Join(categToRem, ","))
	// err = dbAldo.Select(&products, fmt.Sprintf("SELECT * FROM product WHERE category IN (%s)", strings.Join(categToRem, ",")))

	// if err != nil {
	// log.Fatal(fmt.Errorf("Get products to remove from db: %v", err))
	// }

	// for _, p := range products {
	// log.Println("products to remove:", p.Code)
	// }

	// Remove products.
	// dbAldo.MustExec(fmt.Sprintf("DELETE FROM product WHERE category IN (%s)", strings.Join(categToRem, ",")))
}
