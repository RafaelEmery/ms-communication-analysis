package client

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
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
	Resource        string `json:"resource"`
	RequestQuantity int    `json:"request_quantity"`
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
	v1 := a.Group("/interact")
	v1.Post("/http", b.interactWithHTTPServer)
	v1.Post("/grpc", b.interactWithGRPCServer)
	v1.Post("/rabbitmq", b.interactWithRabbitMQ)
}

func (b *BFFApp) interactWithHTTPServer(c *fiber.Ctx) error {
	var data InteractionData
	if err := c.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	o, err := b.handleMethods(c, data, httpMethod)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendString(fmt.Sprintf("[%s] %d req for %s", o, data.RequestQuantity, data.Resource))
}

func (b *BFFApp) interactWithGRPCServer(c *fiber.Ctx) error {
	var data InteractionData
	if err := c.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	o, err := b.handleMethods(c, data, grpcMethod)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendString(fmt.Sprintf("[%s] %d req for %s", o, data.RequestQuantity, data.Resource))
}

func (b *BFFApp) interactWithRabbitMQ(c *fiber.Ctx) error {
	var data InteractionData
	if err := c.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	o, err := b.handleMethods(c, data, rabbitMQMethod)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendString(fmt.Sprintf("[%s] %d req for %s", o, data.RequestQuantity, data.Resource))
}

func (b *BFFApp) handleMethods(c *fiber.Ctx, data InteractionData, method string) (string, error) {
	log.Default().Printf("[%d req] %s method on %s resource\n", data.RequestQuantity, method, data.Resource)

	totalStart := time.Now()
	if method == httpMethod {
		if err := HandleHTTP(b.HTTPBaseURL, data); err != nil {
			return "", err
		}
	}
	if method == grpcMethod {
		if err := HandleGRPC(c.Context(), b.GRPCServerHost, data); err != nil {
			return "", err
		}
	}
	if method == rabbitMQMethod {
		HandleRabbitMQ(b.RabbitMQChannel, b.RabbitMQQueue, data)
	}

	return time.Since(totalStart).String(), nil
}
