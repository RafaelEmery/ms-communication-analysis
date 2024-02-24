package apps

import (
	"context"
	"log"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/gofiber/fiber/v2"
)

type ProductsManager interface {
	Count(ctx context.Context) (int, error)
	BatchCreate(ctx context.Context, p domain.Products) error
	DeleteAll(ctx context.Context) error
}

type ProductsToGenerate struct {
	Quantity int `json:"quantity"`
}

type SetupApp struct {
	ctx context.Context
	d   ProductsManager
}

func NewSetupApp(ctx context.Context, d ProductsManager) SetupApp {
	return SetupApp{ctx: ctx, d: d}
}

func (s *SetupApp) Routes(a *fiber.App) {
	v1 := a.Group("/setup")
	v1.Get("/count", s.getProductsCount)
	v1.Post("/generate", s.generateProducts)
	v1.Delete("/drop", s.dropProducts)
}

func (s *SetupApp) getProductsCount(c *fiber.Ctx) error {
	o, err := s.d.Count(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(o)
}

func (s *SetupApp) generateProducts(c *fiber.Ctx) error {
	var (
		body    ProductsToGenerate
		product domain.Product
	)
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	p := make(domain.Products, 0)
	for i := 0; i < body.Quantity; i++ {
		fp := *product.Fake()
		log.Default().Println("fake product - ", fp)
		p = append(p, fp)
	}

	if err := s.d.BatchCreate(s.ctx, p); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}

func (s *SetupApp) dropProducts(c *fiber.Ctx) error {
	if err := s.d.DeleteAll(context.Background()); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
