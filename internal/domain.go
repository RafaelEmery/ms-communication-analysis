package internal

import (
	"time"

	faker "github.com/brianvoe/gofakeit/v6"
)

type Products []Product

type Product struct {
	ID                string    `json:"id,omitempty"`
	Name              string    `json:"name"`
	SKU               string    `json:"sku"`
	SellerName        string    `json:"seller_name"`
	Price             float64   `json:"price"`
	AvailableQuantity int64     `json:"available_quantity"`
	SalesQuantity     int64     `json:"sales_quantity"`
	Active            bool      `json:"active"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

func (p *Product) Fake() *Product {
	p.ID = faker.UUID()
	p.Name = faker.BeerName()
	p.SKU = faker.UUID()
	p.SellerName = faker.Name()
	p.Price = faker.Price(1.0, 999.0)
	p.AvailableQuantity = faker.Int64()
	p.SalesQuantity = faker.Int64()
	p.Active = faker.Bool()
	p.CreatedAt = faker.Date()
	p.UpdatedAt = faker.Date()

	return p
}
