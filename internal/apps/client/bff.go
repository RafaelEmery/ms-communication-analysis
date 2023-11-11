package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	apps "github.com/RafaelEmery/performance-analysis-server/internal/apps/server"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	o, err := b.handleMethods(c, data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendString(fmt.Sprintf("[%s] %d req for %s | %s", o, data.RequestQuantity, data.Resource, data.CommunicationMethod))
}

func (b *BFFApp) handleMethods(c *fiber.Ctx, data InteractionData) (string, error) {
	log.Default().Printf("\n%s method on %s resource\n\n", data.CommunicationMethod, data.Resource)
	var totalStart time.Time

	if data.CommunicationMethod == httpMethod {
		endpoint, method := b.getRequestData(data.Resource)

		totalStart = time.Now()
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
	}
	if data.CommunicationMethod == grpcMethod {
		conn, err := grpc.Dial(b.GRPCServerHost, grpc.WithInsecure())
		if err != nil {
			return "", err
		}
		defer conn.Close()

		client := apps.NewProductHandlerClient(conn)
		totalStart = time.Now()

		for i := 0; i < data.RequestQuantity; i++ {
			start := time.Now()
			code, err := doProcedureCall(c.Context(), data.Resource, client)
			if err != nil {
				log.Default().Println(err)
				continue
			}

			elapsed := time.Since(start).String()
			log.Default().Printf("[%s] %s - %s", code, strings.ToUpper(string(data.Resource[0]))+data.Resource[1:], elapsed)
		}
	}

	return time.Since(totalStart).String(), nil
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

func doProcedureCall(ctx context.Context, resource string, client apps.ProductHandlerClient) (string, error) {
	var err error

	switch resource {
	case createResource:
		var product domain.Product
		fp := *product.Fake()

		req := &apps.CreateProductRequest{
			Name:              fp.Name,
			Sku:               fp.SKU,
			SellerName:        fp.SellerName,
			Price:             float32(fp.Price),
			AvailableDiscount: float32(fp.AvailableDiscount),
			AvailableQuantity: int32(fp.AvailableQuantity),
			SalesQuantity:     int32(fp.SalesQuantity),
			Active:            fp.Active,
		}

		_, err = client.Create(ctx, req)
	case reportResource:
		req := &apps.EmptyRequest{}
		_, err = client.Report(ctx, req)
	case getByDiscountResource:
		req := &apps.EmptyRequest{}
		_, err = client.GetByDiscount(ctx, req)
	}

	return getProcedureCallCode(ctx), err
}

func getProcedureCallCode(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "not_found"
	}
	statusCodes, ok := md["grpc-status"]
	if !ok {
		return "not_found"
	}
	return statusCodes[0]
}
