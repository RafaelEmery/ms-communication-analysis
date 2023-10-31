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
		INSERT INTO products (id, name, sku, seller_name, price, discount, available_quantity, sales_quantity, active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.SKU, p.SellerName, p.Price, p.AvailableDiscount, p.AvailableQuantity, p.SalesQuantity, p.Active, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Get(ctx context.Context) (Products, error) {
	query := `
		SELECT 
			id, name, sku, seller_name, price, discount, available_quantity, sales_quantity, active, created_at, updated_at
		FROM products
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var p Products
	for rows.Next() {
		var i Product
		if err := rows.Scan(&i.ID, &i.Name, &i.SKU, &i.SellerName, &i.Price, &i.AvailableDiscount, &i.AvailableQuantity, &i.SalesQuantity, &i.Active, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return p, err
		}
		p = append(p, i)
	}
	if err = rows.Err(); err != nil {
		return p, err
	}
	return p, nil
}

func (r Repository) Count(ctx context.Context) (count int, err error) {
	query := `SELECT COUNT(id) FROM products`

	if err = r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return
	}

	return
}

func (r Repository) BatchCreate(ctx context.Context, p Products) error {
	for _, i := range p {
		query := `
			INSERT INTO products (id, name, sku, seller_name, price, discount, available_quantity, sales_quantity, active, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		_, err := r.db.ExecContext(ctx, query, i.ID, i.Name, i.SKU, i.SellerName, i.Price, i.AvailableDiscount, i.AvailableQuantity, i.SalesQuantity, i.Active, i.CreatedAt, i.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Repository) DeleteAll(ctx context.Context) error {
	query := `
		DELETE FROM products
	`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
