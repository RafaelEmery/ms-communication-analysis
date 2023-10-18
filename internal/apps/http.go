package apps

import (
	"context"
	"net/http"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/gofiber/fiber/v2"
)

type Return struct {
	Result any `json:"result"`
}

type HttpApp struct {
	ctx context.Context
	c   Creator
	rg  ReportGenerator
}

func NewHttpApp(ctx context.Context, c Creator, rg ReportGenerator) HttpApp {
	return HttpApp{ctx: ctx, c: c, rg: rg}
}

func (h *HttpApp) Routes(a *fiber.App) {
	v1 := a.Group("/products")
	v1.Post("", h.createProduct)
	v1.Get("/report", h.getReport)
}

func (h *HttpApp) createProduct(c *fiber.Ctx) error {
	var (
		product domain.Product
	)
	if err := c.BodyParser(&product); err != nil {
		return err
	}

	o, err := h.c.Create(h.ctx, &product)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(Return{Result: o})
}

func (h *HttpApp) getReport(c *fiber.Ctx) error {
	o, err := h.rg.GenerateReport(h.ctx)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(Return{Result: o})
}
