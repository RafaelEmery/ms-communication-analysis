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
	AvailableDiscount float64   `json:"available_discount"`
	AvailableQuantity int       `json:"available_quantity"`
	SalesQuantity     int       `json:"sales_quantity"`
	Active            bool      `json:"active"`
	DiscountApplied   bool      `json:"discount_applied"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

func (p *Product) Fake() *Product {
	p.ID = faker.UUID()
	p.Name = faker.BeerName()
	p.SKU = faker.UUID()
	p.SellerName = faker.Name()
	p.Price = faker.Price(1.0, 999.0)
	p.AvailableDiscount = faker.Float64Range(0.0, 1.0)
	p.AvailableQuantity = faker.RandomInt([]int{100, 200, 300, 400, 500})
	p.SalesQuantity = faker.RandomInt([]int{10, 20, 30, 40, 50})
	p.Active = faker.Bool()
	p.DiscountApplied = false
	p.CreatedAt = faker.Date()
	p.UpdatedAt = faker.Date()

	return p
}
