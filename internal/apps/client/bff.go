package client

import (
	"fmt"
	"log"
	"runtime"
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

type InteractionData struct {
	Resource        string `json:"resource"`
	RequestQuantity int    `json:"request_quantity"`
}

type InteractionInfo struct {
	MemoryUsage uint64
	RequestTime string
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
	a.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendString("OK!")
	})
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

	output := fmt.Sprintf("request time - %s | memory - %d (kB)", o.RequestTime, o.MemoryUsage)
	log.Println(output)
	return c.SendString(output)
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

	output := fmt.Sprintf("request time - %s | memory - %d (kB)", o.RequestTime, o.MemoryUsage)
	log.Println(output)
	return c.SendString(output)
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

	output := fmt.Sprintf("request time - %s | memory - %d (kB)", o.RequestTime, o.MemoryUsage)
	log.Println(output)
	return c.SendString(output)
}

func (b *BFFApp) handleMethods(c *fiber.Ctx, data InteractionData, method string) (InteractionInfo, error) {
	totalStart := time.Now()

	if method == httpMethod {
		if err := HandleHTTP(b.HTTPBaseURL, data); err != nil {
			return InteractionInfo{}, err
		}
	}
	if method == grpcMethod {
		if err := HandleGRPC(c.Context(), b.GRPCServerHost, data); err != nil {
			return InteractionInfo{}, err
		}
	}
	if method == rabbitMQMethod {
		HandleRabbitMQ(b.RabbitMQChannel, b.RabbitMQQueue, data)
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	info := InteractionInfo{
		RequestTime: time.Since(totalStart).String(),
		MemoryUsage: convertToKB(memStats.TotalAlloc),
	}

	return info, nil
}

func convertToKB(v uint64) uint64 {
	return v / 1024.0
}
