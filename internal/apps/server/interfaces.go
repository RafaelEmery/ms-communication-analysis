package apps

import (
	"context"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

type Creator interface {
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
}

type ReportGenerator interface {
	GenerateReport(ctx context.Context) (string, error)
}

type ProductByDiscountGetter interface {
	GetByDiscount(ctx context.Context) (domain.Products, error)
}
