package usecases

import (
	"context"
	"fmt"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
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

func (u reportUseCase) GenerateReport(ctx context.Context) (string, error) {
	p, err := u.g.Get(ctx)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("./.tmp/report_%s.pdf", time.Now().Format("2006-02-01_18:00:00"))
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	for _, product := range p {
		m.Row(5, func() {
			m.Col(12, func() {
				text := fmt.Sprintf("Product: %s - %s - %f - %d - %d", product.Name, product.SKU, product.Price, product.SalesQuantity, product.AvailableQuantity)
				m.Text(text, props.Text{
					Size:            10,
					Align:           consts.Center,
					Top:             50,
					VerticalPadding: 2.0,
				})
			})
		})
	}

	if err := m.OutputFileAndClose(fileName); err != nil {
		return "", err
	}

	return fileName, nil
}
