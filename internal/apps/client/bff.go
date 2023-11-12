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
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	httpMethod            = "http"
	grpcMethod            = "grpc"
	rabbitMQMethod        = "rabbitmq"
	createResource        = "create"
	reportResource        = "report"
	getByDiscountResource = "getByDiscount"
)

type Message struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type InteractionData struct {
	Resource            string `json:"resource"`
	CommunicationMethod string `json:"communication_method"`
	RequestQuantity     int    `json:"request_quantity"`
}

type BFFApp struct {
	HTTPBaseURL     string
	GRPCServerHost  string
	RabbitMQChannel *amqp.Channel
	RabbitMQQueue   amqp.Queue
}

func NewBFFApp(h, g string, ch *amqp.Channel, q amqp.Queue) BFFApp {
	return BFFApp{HTTPBaseURL: h, GRPCServerHost: g, RabbitMQChannel: ch, RabbitMQQueue: q}
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
		endpoint, method := getRequestData(data.Resource, b.HTTPBaseURL)

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
	if data.CommunicationMethod == rabbitMQMethod {
		for i := 0; i < data.RequestQuantity; i++ {
			jsonBody, err := getMessageBody(data.Resource)
			if err != nil {
				log.Default().Println(err)
				continue
			}
			if err := b.RabbitMQChannel.Publish(
				"",                   // Exchange
				b.RabbitMQQueue.Name, // Routing key (nome da fila)
				false,                // Mandatory
				false,                // Immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        jsonBody,
					Headers: amqp.Table{
						"resource": data.Resource,
					},
				}); err != nil {
				log.Default().Println(err)
				continue
			}

		}
	}

	return time.Since(totalStart).String(), nil
}

func getRequestData(resource, baseURL string) (string, string) {
	var (
		endpoint      string
		requestMethod string
	)
	switch resource {
	case createResource:
		endpoint = baseURL + "/products"
		requestMethod = http.MethodPost
	case reportResource:
		endpoint = baseURL + "/products/report"
		requestMethod = http.MethodGet
	case getByDiscountResource:
		endpoint = baseURL + "/products/discount"
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

// TODO: procedure call code is not working
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

func getMessageBody(resource string) ([]byte, error) {
	strContent := ""
	if resource == createResource {
		var product domain.Product

		content, err := json.Marshal(product.Fake())
		if err != nil {
			return []byte{}, err
		}

		strContent = string(content)
	}

	m := Message{
		Content: strContent,
		Metadata: map[string]string{
			"resource": resource,
		},
	}

	body, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
