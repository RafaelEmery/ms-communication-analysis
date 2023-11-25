package usecases

import (
	"context"
	"log"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/google/uuid"
)

type Creator interface {
	Create(ctx context.Context, p domain.Product) error
}

type createUseCase struct {
	c Creator
}

func NewCreateUseCase(c Creator) createUseCase {
	return createUseCase{c: c}
}

func (u createUseCase) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	p.ID = uuid.NewString()
	p.DiscountApplied = false
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	start := time.Now()
	if err := u.c.Create(ctx, p); err != nil {
		return domain.Product{}, err
	}
	log.Printf("database interaction time - %s", time.Since(start).String())

	return p, nil
}
