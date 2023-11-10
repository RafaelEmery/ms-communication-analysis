package usecases

import (
	"context"
	"fmt"
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

// TODO: should product be a pointer?
func (u createUseCase) Create(ctx context.Context, p domain.Product) (domain.Product, error) {

	fmt.Println("inside use case")

	p.ID = uuid.NewString()
	p.DiscountApplied = false
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if err := u.c.Create(ctx, p); err != nil {
		fmt.Println("error on creating", err.Error())
		return domain.Product{}, err
	}

	fmt.Println("product - ", p)

	return p, nil
}
