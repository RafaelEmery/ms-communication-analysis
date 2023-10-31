package usecases

import (
	"context"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

type ProductGetter interface {
	Get(ctx context.Context) (domain.Products, error)
}
