package main

import (
	"fmt"
	"log"

	goenv "github.com/Netflix/go-env"
	"github.com/RafaelEmery/performance-analysis-server/internal/apps/client"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // Load .env file
	"github.com/streadway/amqp"
)

type Env struct {
	ClientPort string `env:"BFF_APP_PORT"`
	HTTPPort   string `env:"HTTP_APP_PORT"`
	GRPCPort   string `env:"GRPC_APP_PORT"`
	RabbitMQ   struct {
		Port      string `env:"RABBITMQ_PORT"`
		WebPort   string `env:"RABBITMQ_WEB_PORT"`
		User      string `env:"RABBITMQ_USER"`
		QueueName string `env:"RABBITMQ_QUEUE_NAME"`
		Host      string `env:"RABBITMQ_HOST"`
	}
}

func getEnv() (*Env, error) {
	env := &Env{}
	_, err := goenv.UnmarshalFromEnviron(env)

	log.Default().Println("environment loaded - ", env)

	return env, err
}

func main() {
	env, err := getEnv()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	httpURL := fmt.Sprintf("http://localhost:%s", env.HTTPPort)
	grpcHost := fmt.Sprintf("localhost:%s", env.GRPCPort)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", env.RabbitMQ.User, env.RabbitMQ.User, env.RabbitMQ.Host, env.RabbitMQ.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		env.RabbitMQ.QueueName, // Name
		false,                  // Durable
		false,                  // Delete when unused
		false,                  // Exclusive
		false,                  // No-wait
		nil,                    // Arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	bff := client.NewBFFApp(httpURL, grpcHost, ch, q)
	bff.Routes(app)

	log.Default().Println("client application working")
	app.Listen(fmt.Sprintf(":%s", env.ClientPort))
}
