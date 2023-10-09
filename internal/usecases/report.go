package usecases

import (
	"context"
	"fmt"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

type SalesQuantityGetter interface {
	Get(ctx context.Context) (domain.Products, error)
}

type reportUseCase struct {
	g SalesQuantityGetter
}

func NewReportUseCase(g SalesQuantityGetter) reportUseCase {
	return reportUseCase{g: g}
}

func (u reportUseCase) GenerateReport(ctx context.Context) ([]byte, error) {
	p, err := u.g.Get(ctx)
	if err != nil {
		return []byte{}, err
	}

	// TODO: implement file generation with three fields
	for _, v := range p {
		fmt.Println(v.Name)
		fmt.Println(v.SKU)
		fmt.Println(v.SalesQuantity)
	}

	return []byte{}, nil
}
