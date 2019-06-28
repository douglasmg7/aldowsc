package main

import (
	"bytes"
	"reflect"
	"strings"
	"time"

	"github.com/douglasmg7/money"
)

type Product struct {
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

func (p *Product) Find(Id string) error {
	var fieldsName []string
	var fieldsNameDb []string
	var fieldsInterface []interface{}

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
	buffer.WriteString("product ")
	buffer.WriteString("WHERE code=?")

	err := db.QueryRow(buffer.String(), Id).Scan(fieldsInterface...)
	return err
}

func (p *Product) Save() error {
	var fieldsName []string
	var fieldsNameDb []string
	var fieldsInterface []interface{}

	val := reflect.ValueOf(p).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		fieldsName = append(fieldsName, fieldType.Name)
		fieldsNameDb = append(fieldsNameDb, fieldType.Tag.Get("db"))
		fieldsInterface = append(fieldsInterface, val.Field(i).Addr().Interface())
	}
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO ")
	buffer.WriteString("product ")
	buffer.WriteString(`(`)
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
	if err != nil {
		return err
	}
	return err
}

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
