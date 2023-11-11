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
	c  Creator
	rg ReportGenerator
	pg ProductByDiscountGetter
}

func NewHttpApp(ctx context.Context, c Creator, rg ReportGenerator, pg ProductByDiscountGetter) HttpApp {
	return HttpApp{c: c, rg: rg, pg: pg}
}

func (h *HttpApp) Routes(a *fiber.App) {
	v1 := a.Group("/products")
	v1.Post("", h.createProduct)
	v1.Get("/report", h.getReport)
	v1.Get("/discount", h.getByAppliedDiscount)
}

func (h *HttpApp) createProduct(c *fiber.Ctx) error {
	var (
		product domain.Product
	)
	if err := c.BodyParser(&product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	o, err := h.c.Create(c.Context(), product)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(Return{Result: o})
}

func (h *HttpApp) getReport(c *fiber.Ctx) error {
	o, err := h.rg.GenerateReport(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(Return{Result: o})
}

// TODO: test get by applied discount endpoint
func (h *HttpApp) getByAppliedDiscount(c *fiber.Ctx) error {
	o, err := h.pg.GetByDiscount(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(Return{Result: o})
}
