package usecases

import (
	"context"
	"sort"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

type getByDiscountUseCase struct {
	g ProductGetter
}

func NewGetByDiscountUseCase(g ProductGetter) getByDiscountUseCase {
	return getByDiscountUseCase{g: g}
}

func (u getByDiscountUseCase) GetByDiscount(ctx context.Context) (domain.Products, error) {
	p, err := u.g.Get(ctx)
	if err != nil {
		return domain.Products{}, err
	}

	dps := u.applyAvailableDiscount(p)
	orderedDPs := u.orderByPrice(dps)

	return orderedDPs, nil
}

func (u getByDiscountUseCase) applyAvailableDiscount(ps domain.Products) domain.Products {
	dp := make(domain.Products, 0)
	for _, p := range ps {
		discountPrice := p.Price - (p.Price / p.AvailableDiscount)
		p.Price = discountPrice
		p.DiscountApplied = true

		dp = append(dp, p)
	}

	return dp
}

func (u getByDiscountUseCase) orderByPrice(dps domain.Products) domain.Products {
	sort.Slice(dps, func(i, j int) bool {
		return dps[i].Price < dps[j].Price
	})

	return dps
}
