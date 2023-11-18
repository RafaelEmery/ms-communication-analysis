package main

import (
	"fmt"
	"log"
	"time"

	goenv "github.com/Netflix/go-env"
	"github.com/RafaelEmery/performance-analysis-server/internal/apps/client"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // Load .env file
	"github.com/streadway/amqp"
)

type Env struct {
	HTTPHost   string `env:"HTTP_HOST"`
	GRPCHost   string `env:"GRPC_HOST"`
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
	httpURL := fmt.Sprintf("http://%s:%s", env.HTTPHost, env.HTTPPort)
	grpcHost := fmt.Sprintf("%s:%s", env.GRPCHost, env.GRPCPort)

	time.Sleep(5 * time.Second)
	var conn *amqp.Connection
	var ch *amqp.Channel
	var q amqp.Queue

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", env.RabbitMQ.User, env.RabbitMQ.User, env.RabbitMQ.Host, env.RabbitMQ.Port)
	log.Default().Println("rabbitMQ connection string: ", connectionString)
	conn, err = amqp.Dial(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Default().Println("rabbitMQ successfully connected: ", !conn.IsClosed())

	ch, err = conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err = ch.QueueDeclare(
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

	// TODO: pass all connections and hosts to BFF app and handle methods/calls only

	bff := client.NewBFFApp(httpURL, grpcHost, ch, q)
	bff.Routes(app)

	log.Default().Println("client application working")
	app.Listen(fmt.Sprintf(":%s", env.ClientPort))
}
