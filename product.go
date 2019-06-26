package main

import (
	"time"

	"github.com/douglasmg7/money"
)

type Product struct {
	Code                 string
	Brand                string
	Category             string
	Description          string
	Multiple             int
	DealerPrice          money.Money
	SuggestionPrice      money.Money
	TechnicalDescription string
	Availability         bool
	Length               int // mm.
	Width                int // mm.
	Height               int // mm.
	Weight               int // grams.
	PictureLink          string
	WarrantyPeriod       int    // Days.
	RMAProcedure         string // ?
	CreatedAt            time.Time
	ChangedAt            time.Time
	Changed              bool
	New                  bool
	Removed              bool
}

func (p *Product) Find(Id string) error {
	err := db.QueryRow(`
		SELECT 
			code, 
			brand, 
			category, 
			description, 
			multiple,
			dealer_price,
			suggestion_price,
			technical_description,
			availability, 
			length,
			width,
			height,
			weight,
			picture_link,
			warranty_period,
			rma_procedure,
			created_at,
			changed_at,
			changed,
			new,
			removed
		FROM 
			product 
		WHERE 
			code = ?`, Id).
		Scan(
			&p.Code,
			&p.Brand,
			&p.Category,
			&p.Description,
			&p.Multiple,
			&p.DealerPrice,
			&p.SuggestionPrice,
			&p.TechnicalDescription,
			&p.Availability,
			&p.Length,
			&p.Width,
			&p.Height,
			&p.Weight,
			&p.PictureLink,
			&p.WarrantyPeriod,
			&p.RMAProcedure,
			&p.CreatedAt,
			&p.ChangedAt,
			&p.Changed,
			&p.New,
			&p.Removed)
	return err
}

func (p *Product) Save() error {
	now := time.Now()
	stmt, err := db.Prepare(`
		INSERT INTO product(
			code, 
			brand, 
			category, 
			description, 
			multiple,
			dealer_price,
			suggestion_price,
			technical_description,
			availability, 
			length,
			width,
			height,
			weight,
			picture_link,
			warranty_period,
			rma_procedure,
			created_at,
			changed_at,
			changed,
			new,
			removed
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		p.Code,
		p.Brand,
		p.Category,
		p.Description,
		p.Multiple,
		p.DealerPrice,
		p.SuggestionPrice,
		p.TechnicalDescription,
		p.Availability,
		p.Length,
		p.Width,
		p.Height,
		p.Weight,
		p.PictureLink,
		p.WarrantyPeriod,
		p.RMAProcedure,
		now,
		now,
		false,
		true,
		false)
	if err != nil {
		return err
	}
	return err
}
