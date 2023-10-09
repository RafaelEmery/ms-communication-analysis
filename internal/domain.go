package internal

import "time"

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

type Products []Product
