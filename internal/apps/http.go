package apps

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

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
	v1.Get("/report", h.getReport)
	v1.Post("", h.createProduct)
}

func (h *HttpApp) createProduct(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(nil)
}

func (h *HttpApp) getReport(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(nil)
}
