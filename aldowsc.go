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
	"unicode"

	"github.com/douglasmg7/currency"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/html/charset"
)

var db *sql.DB
var dbAldo *sqlx.DB

// Paths.
var appPath, logPath, dbPath, xmlPath, listPath string

// Files.
var logFile, dbAldoFile string

// Min and max price filter.
var maxPriceFilter, minPriceFilter currency.Currency

// Development mode.
var dev bool

// Categories selected to use.
var categSel []string

func init() {
	// Path.
	zunkaPath := os.Getenv("ZUNKAPATH")
	if zunkaPath == "" {
		panic("ZUNKAPATH not defined.")
	}
	logPath := path.Join(zunkaPath, "log")
	listPath = path.Join(zunkaPath, "list")
	xmlPath = path.Join(zunkaPath, "xml")
	// Create path.
	os.MkdirAll(logPath, os.ModePerm)
	os.MkdirAll(listPath, os.ModePerm)
	os.MkdirAll(xmlPath, os.ModePerm)

	// Aldo db.
	dbAldoFile = os.Getenv("ZUNKA_ALDOWSC_DB")
	if dbAldoFile == "" {
		panic("ZUNKA_ALDOWSC_DB not defined.")
	}
	dbAldoFile = path.Join(zunkaPath, "db", dbAldoFile)

	// Log file.
	logFile, err := os.OpenFile(path.Join(logPath, "aldowsc.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// Filters.
	minPriceFilter, err = currency.Parse("2.000,00", ",")
	if err != nil {
		panic(err)
	}
	maxPriceFilter, err = currency.Parse("100.000,00", ",")
	if err != nil {
		panic(err)
	}

	// Log configuration.
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Ldate | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Development mode.
	mode := ""
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "dev") {
		dev = true
		mode = "development"
	} else {
		mode = "production"
	}

	// Log start.
	log.Printf("*** Starting aldowsc in %v mode (version %s) ***\n", mode, version)
}

func main() {
	var err error

	// Db start.
	db, err = sql.Open("sqlite3", dbAldoFile)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Init db Aldo.
	dbAldo = sqlx.MustConnect("sqlite3", dbAldoFile)

	// Read selected categories.
	log.Println("Reading selected categories list...")
	categSel = readList(path.Join(listPath, "categSel.list"))
	// log.Println("categSel:", categSel)

	// Remove no more selected products.
	rmProductsNotSel()

	// Remove products with price out of range.
	rmProductsPriceOutOfRange()

	// Load xml file.
	log.Println("Loading and decoding xml file...")
	// timer := time.Now()
	aldoXMLDoc := xmlDoc{}
	decoder := xml.NewDecoder(os.Stdin)

	// Decoding xml file.
	timer := time.Now()
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&aldoXMLDoc)
	if err != nil {
		log.Fatalln("Error decoding xml file:", err)
	}
	log.Printf("Time loading and decoding xml file: %fs", time.Since(timer).Seconds())

	// Not process zero products, to not remove current product.
	if len(aldoXMLDoc.Products) == 0 {
		log.Println("*** XML returned zero products ***")
		log.Printf("Finish\n\n")
		return
	}

	// Processing products.
	log.Println("Processing products...")
	timer = time.Now()
	err = aldoXMLDoc.process()
	log.Printf("Time processing products: %fs", time.Since(timer).Seconds())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Finish\n\n")
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

	// Each no blank line is um item without spaces.
	sa := []string{}
	ra := []rune{}
	for _, r := range string(b) {
		if r == '\n' {
			if len(ra) > 0 {
				sa = append(sa, string(ra))
				ra = ra[:0]
			}
			// fmt.Println("new line")
		} else if !unicode.IsSpace(r) {
			ra = append(ra, r)
			// fmt.Printf("\"%c\" %v index %v\n", r, r, index)
		}
	}
	if len(ra) > 0 {
		sa = append(sa, string(ra))
		ra = nil
	}
	// fmt.Println("sa: ", sa)

	return sa
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
		// log.Printf("list: %s, xml: %s", categItem, categorie)
		if strings.HasPrefix(categorie, categItem) {
			// log.Println("selected")
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
		log.Fatal(fmt.Errorf("Get distinct categories from db. %v", err))
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
}

// Remove products with price out of defined range.
func rmProductsPriceOutOfRange() {
	// Copy products to remove to history.
	tx := dbAldo.MustBegin()
	// tx.MustExec(fmt.Sprintf("INSERT INTO product_history SELECT * FROM product WHERE category IN (%s)", strings.Join(categToRem, ",")))
	tx.MustExec(fmt.Sprintf("INSERT INTO product_history SELECT * FROM product WHERE dealer_price NOT BETWEEN (%d) AND (%d)", minPriceFilter, maxPriceFilter))
	// Delete copied products.
	result := tx.MustExec(fmt.Sprintf("DELETE FROM product WHERE dealer_price NOT BETWEEN (%d) AND (%d)", minPriceFilter, maxPriceFilter))
	err := tx.Commit()
	if err != nil {
		log.Fatal(fmt.Errorf("Removing products from db: %v", err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(fmt.Errorf("Removing products from db: %v", err))
	}
	if rowsAffected > 0 {
		log.Printf("Removed %v product(s) with price out of range", rowsAffected)
	}
}
