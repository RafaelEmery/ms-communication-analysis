package apps

import (
	"context"
	"log"
	"runtime"

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

type InteractionInfo struct {
	MemoryUsage uint64
}

func logMemStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	i := InteractionInfo{MemoryUsage: convertToKB(memStats.TotalAlloc)}
	log.Default().Printf("memory - %d (kB)", i.MemoryUsage)
}

func convertToKB(v uint64) uint64 {
	return v / 1024.0
}
