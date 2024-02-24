# Performance Analysis Applications

Base project 

Just run the following query for `products` table:

```sql
CREATE TABLE IF NOT EXISTS
  products (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    sku VARCHAR(64) NOT NULL,
    seller_name VARCHAR(64) NOT NULL,
    price FLOAT NOT NULL,
    available_discount FLOAT NOT NULL,
    available_quantity INTEGER NOT NULL,
    sales_quantity INTEGER NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL default now(),
    updated_at TIMESTAMPTZ NOT NULL default now(),
    CONSTRAINT product_pkey PRIMARY KEY (id)
);
```