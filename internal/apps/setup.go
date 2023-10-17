package apps

import (
	"context"
	"net/http"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/gofiber/fiber/v2"
)

type ProductsManager interface {
	Count(ctx context.Context) (int, error)
	BatchCreate(ctx context.Context, p domain.Products) error
	DeleteAll(ctx context.Context) error
}

type ProductsToGenerate struct {
	Quantity int
}

type ProjectContext struct {
	ctx context.Context
	d   ProductsManager
}

func (pc *ProjectContext) Routes(a *fiber.App) {
	a.Group("/setup")
	a.Get("/quantity", pc.getProductsQuantity)
	a.Post("/generate", pc.generateProducts)
	a.Post("/drop", pc.dropProducts)
}

func (pc *ProjectContext) getProductsQuantity(c *fiber.Ctx) error {
	o, err := pc.d.Count(context.Background())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(o)
}

func (pc *ProjectContext) generateProducts(c *fiber.Ctx) error {
	var (
		body    ProductsToGenerate
		product domain.Product
	)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	p := make(domain.Products, 0)
	for i := 0; i < body.Quantity; i++ {
		p = append(p, *product.Fake())
	}

	if err := pc.d.BatchCreate(pc.ctx, p); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusNoContent).JSON(nil)
}

func (pc *ProjectContext) dropProducts(c *fiber.Ctx) error {
	if err := pc.d.DeleteAll(context.Background()); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusNoContent).JSON(nil)
}
