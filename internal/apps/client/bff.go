package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"github.com/gofiber/fiber/v2"
)

const (
	httpMethod            = "http"
	grpcMethod            = "grpc"
	createResource        = "create"
	reportResource        = "report"
	getByDiscountResource = "getByDiscount"
)

type InteractionData struct {
	Resource            string `json:"resource"`
	CommunicationMethod string `json:"communication_method"`
	RequestQuantity     int    `json:"request_quantity"`
}

type BFFApp struct {
	HTTPBaseURL    string
	GRPCServerHost string
}

func NewBFFApp(h string, g string) BFFApp {
	return BFFApp{HTTPBaseURL: h, GRPCServerHost: g}
}

func (b *BFFApp) Routes(a *fiber.App) {
	a.Post("/interact", b.interactWithServer)
}

func (b *BFFApp) interactWithServer(c *fiber.Ctx) error {
	var data InteractionData
	if err := c.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return b.handleMethods(c, data)
}

func (b *BFFApp) handleMethods(c *fiber.Ctx, data InteractionData) error {
	if data.CommunicationMethod == httpMethod {
		endpoint, method := b.getRequestData(data.Resource)

		for i := 0; i < data.RequestQuantity; i++ {
			start := time.Now()
			resp, err := doRequest(endpoint, method)
			if err != nil {
				log.Default().Println(err)
				continue
			}
			defer resp.Body.Close()

			elapsed := time.Since(start).String()
			log.Default().Printf("[%d] %s - %s", resp.StatusCode, endpoint, elapsed)
		}

		return c.SendStatus(fiber.StatusOK)
	}
	if data.CommunicationMethod == grpcMethod {
		return nil
	}

	return nil
}

func (b *BFFApp) getRequestData(resource string) (string, string) {
	var (
		endpoint      string
		requestMethod string
	)
	switch resource {
	case createResource:
		endpoint = b.HTTPBaseURL + "/products"
		requestMethod = http.MethodPost
	case reportResource:
		endpoint = b.HTTPBaseURL + "/products/report"
		requestMethod = http.MethodGet
	case getByDiscountResource:
		endpoint = b.HTTPBaseURL + "/products/discount"
		requestMethod = http.MethodGet
	}

	return endpoint, requestMethod
}

func doRequest(endpoint, method string) (*http.Response, error) {
	if method == http.MethodPost {
		var product domain.Product

		payload, err := json.Marshal(product.Fake())
		if err != nil {
			log.Default().Println(err)
			return nil, err
		}
		body := bytes.NewBuffer(payload)

		return http.Post(endpoint, "application/json", body)
	}
	if method == http.MethodGet {
		return http.Get(endpoint)
	}

	return nil, fmt.Errorf("the method %s is not allowed", method)
}
