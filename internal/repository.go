package internal

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(ctx context.Context, p Product) error {
	query := `
		INSERT INTO products (id, name, sku, seller_name, price, available_quantity, sales_quantity, active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.SKU, p.SellerName, p.Price, p.AvailableQuantity, p.SalesQuantity, p.Active, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Get(ctx context.Context) (Products, error) {
	return Products{}, nil
}
