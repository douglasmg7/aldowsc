package main

import (
	"time"

	"github.com/douglasmg7/money"
)

type Product struct {
	Code                string
	Brand               string
	Category            string
	Description         string
	Multiple            int
	DealerPrice         money.Money
	SuggestionPrice     money.Money
	TecnicalDescription string
	Availability        bool
	Length              int // mm.
	Width               int // mm.
	Height              int // mm.
	Weight              int // grams.
	PictureLinks        string
	WarrantyPeriod      int    // Days.
	RMAProcedure        string // ?
	CreatedAt           time.Time
	ChangedAt           time.Time
	Changed             bool
	New                 bool
	Removed             bool
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
			tecnical_description,
			availability, 
			length,
			width,
			height,
			weight,
			picture_links,
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
			&p.TecnicalDescription,
			&p.Availability,
			&p.Length,
			&p.Width,
			&p.Height,
			&p.Weight,
			&p.PictureLinks,
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
	stmt, err := db.Prepare(`
		INSERT INTO product(
			code, 
			brand, 
			category, 
			description, 
			multiple,
			dealer_price,
			suggestion_price,
			tecnical_description,
			availability, 
			length,
			width,
			height,
			weight,
			picture_links,
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
		p.TecnicalDescription,
		p.Availability,
		p.Length,
		p.Width,
		p.Height,
		p.Weight,
		p.PictureLinks,
		p.WarrantyPeriod,
		p.RMAProcedure,
		time.Now(),
		time.Now(),
		false,
		true,
		false)
	if err != nil {
		return err
	}
	return err
}
