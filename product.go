package main

import (
	"bytes"
	"reflect"
	"strings"
	"time"

	"github.com/douglasmg7/money"
)

type Product struct {
	Id                   int         `db:"id"`
	Code                 string      `db:"code"`
	Brand                string      `db:"brand"`
	Category             string      `db:"category"`
	Description          string      `db:"description"`
	Unit                 string      `db:"unit"`
	Multiple             int         `db:"multiple"`
	DealerPrice          money.Money `db:"dealer_price"`
	SuggestionPrice      money.Money `db:"suggestion_price"`
	TechnicalDescription string      `db:"technical_description"`
	Availability         bool        `db:"availability"`
	Length               int         `db:"length"` // mm.
	Width                int         `db:"width"`  // mm.
	Height               int         `db:"height"` // mm.
	Weight               int         `db:"weight"` // grams.
	PictureLink          string      `db:"picture_link"`
	WarrantyPeriod       int         `db:"warranty_period"` // Days.
	RMAProcedure         string      `db:"rma_procedure"`
	CreatedAt            time.Time   `db:"created_at"`
	ChangedAt            time.Time   `db:"changed_at"`
	Changed              bool        `db:"changed"`
	New                  bool        `db:"new"`
	Removed              bool        `db:"removed"`
}

// FindByCode get product from db by code.
func (p *Product) FindByCode(code string) error {
	return p.findByCode(code, false)
}

// FindHistoryByCode get product history from db by code.
func (p *Product) FindHistoryByCode(code string) error {
	return p.findByCode(code, true)
}

// findByCode get product or product history from db by code.
func (p *Product) findByCode(code string, history bool) error {
	var fieldsName []string
	var fieldsNameDb []string
	var fieldsInterface []interface{}
	// Table name.
	var tableName = "product"
	if history {
		tableName = "product_history"
	}

	val := reflect.ValueOf(p).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		fieldsName = append(fieldsName, fieldType.Name)
		fieldsNameDb = append(fieldsNameDb, fieldType.Tag.Get("db"))
		fieldsInterface = append(fieldsInterface, val.Field(i).Addr().Interface())
	}
	var buffer bytes.Buffer
	buffer.WriteString("SELECT ")
	buffer.WriteString(strings.Join(fieldsNameDb, ", "))
	buffer.WriteString(" FROM ")
	buffer.WriteString(tableName)
	buffer.WriteString(" WHERE code=?")

	err := db.QueryRow(buffer.String(), code).Scan(fieldsInterface...)
	return err
}

// Save product to db.
func (p *Product) Save() error {
	return p.save(false)
}

// Save product history to db.
func (p *Product) SaveHistory() error {
	return p.save(true)
}

// Save  product or product history to db.
func (p *Product) save(history bool) error {
	var fieldsName []string
	var fieldsNameDb []string
	var fieldsInterface []interface{}
	// Table name.
	var tableName = "product"
	if history {
		tableName = "product_history"
	}

	val := reflect.ValueOf(p).Elem()
	// i=1, let db generate id.
	for i := 1; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		fieldsName = append(fieldsName, fieldType.Name)
		fieldsNameDb = append(fieldsNameDb, fieldType.Tag.Get("db"))
		fieldsInterface = append(fieldsInterface, val.Field(i).Addr().Interface())
	}
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(tableName)
	buffer.WriteString(` (`)
	buffer.WriteString(strings.Join(fieldsNameDb, ", "))
	buffer.WriteString(`) VALUES(?`)
	buffer.WriteString(strings.Repeat(`, ?`, len(fieldsNameDb)-1))
	buffer.WriteString(`)`)

	stmt, err := db.Prepare(buffer.String())
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(fieldsInterface...)
	return err
}

// Update product from db.
func (p *Product) Update() error {
	return p.update(false)
}

// UpdateHistory update product history db.
func (p *Product) UpdateHistory() error {
	return p.update(true)
}

// Update product or product history db.
func (p *Product) update(history bool) error {
	var fieldsNameSet []string
	// var fieldsNameDb []string
	var fieldsInterface []interface{}
	// Table name.
	var tableName = "product"
	if history {
		tableName = "product_history"
	}

	val := reflect.ValueOf(p).Elem()
	// i=1, to not update id.
	for i := 1; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		// fieldsName = append(fieldsName, fieldType.Name)
		fieldsNameSet = append(fieldsNameSet, fieldType.Tag.Get("db")+"=?")
		fieldsInterface = append(fieldsInterface, val.Field(i).Addr().Interface())
	}
	fieldsInterface = append(fieldsInterface, p.Id)

	query := "UPDATE " + tableName + " SET " + strings.Join(fieldsNameSet, ", ") + " WHERE id=?"
	// fmt.Println(query)
	// fmt.Println("brand:", p.Brand)

	_, err := db.Exec(query, fieldsInterface...)
	return err
}

// Diff check if products are different.
func (p *Product) Diff(pn *Product) bool {
	if p.Code != pn.Code {
		return true
	}
	if p.Brand != pn.Brand {
		return true
	}
	if p.Category != pn.Category {
		return true
	}
	if p.Description != pn.Description {
		return true
	}
	if p.Unit != pn.Unit {
		return true
	}
	if p.Multiple != pn.Multiple {
		return true
	}
	if p.DealerPrice != pn.DealerPrice {
		return true
	}
	if p.SuggestionPrice != pn.SuggestionPrice {
		return true
	}
	if p.TechnicalDescription != pn.TechnicalDescription {
		return true
	}
	if p.Availability != pn.Availability {
		return true
	}
	if p.Length != pn.Length {
		return true
	}
	if p.Width != pn.Width {
		return true
	}
	if p.Height != pn.Height {
		return true
	}
	if p.Weight != pn.Weight {
		return true
	}
	if p.PictureLink != pn.PictureLink {
		return true
	}
	if p.WarrantyPeriod != pn.WarrantyPeriod {
		return true
	}
	if p.RMAProcedure != pn.RMAProcedure {
		return true
	}
	return false
}
