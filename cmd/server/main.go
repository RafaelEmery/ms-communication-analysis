package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	goenv "github.com/Netflix/go-env"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"

	i "github.com/RafaelEmery/performance-analysis-server/internal"
	a "github.com/RafaelEmery/performance-analysis-server/internal/apps/server"
	u "github.com/RafaelEmery/performance-analysis-server/internal/usecases"
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // Load .env file
	_ "github.com/lib/pq"                 // Load postgres connection
)

const (
	setupFlag    = "setup"
	httpFlag     = "http"
	grpcFlag     = "grpc"
	rabbitMQFlag = "rabbitmq"
)

type Env struct {
	AppPorts struct {
		Setup string `env:"SETUP_APP_PORT"`
		HTTP  string `env:"HTTP_APP_PORT"`
		GRPC  string `env:"GRPC_APP_PORT"`
	}
	DB struct {
		Driver   string `env:"DB_DRIVER"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_DATABASE"`
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
	}
	RabbitMQ struct {
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

func connectDatabase(env *Env) (*sql.DB, error) {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DB.Host, env.DB.Port, env.DB.User, env.DB.Password, env.DB.Name)

	db, err := sql.Open(env.DB.Driver, dbSource)
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Default().Println("database connected")

	return db, nil
}

func main() {
	env, err := getEnv()
	if err != nil {
		log.Fatal(err)
	}

	var methodFlag string
	flag.StringVar(&methodFlag, "type", "", "the type of server to run (setup|http|grpc|rabbitmq)")
	flag.Parse()
	log.Default().Printf("running app by flag - %s", methodFlag)

	db, err := connectDatabase(env)
	if err != nil {
		log.Fatal(err)
	}

	r := i.NewRepository(db)

	c := u.NewCreateUseCase(r)
	rg := u.NewReportUseCase(r)
	dpg := u.NewGetByDiscountUseCase(r)

	if methodFlag == setupFlag {
		app := fiber.New()
		setupApp := a.NewSetupApp(context.Background(), r)
		setupApp.Routes(app)

		log.Default().Println("setup application working")
		app.Listen(fmt.Sprintf(":%s", env.AppPorts.Setup))
	}
	if methodFlag == grpcFlag {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.AppPorts.GRPC))
		if err != nil {
			log.Fatal(err)
		}

		s := grpc.NewServer()
		productHandlerServer := a.NewGRPCServer(c, rg, dpg)

		a.RegisterProductHandlerServer(s, productHandlerServer)

		log.Printf("grpc server listening at %v", lis.Addr())
		if err = s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}
	if methodFlag == httpFlag {
		app := fiber.New()
		httpApp := a.NewHttpApp(context.Background(), c, rg, dpg)
		httpApp.Routes(app)

		log.Default().Println("HTTP endpoints working")
		app.Listen(fmt.Sprintf(":%s", env.AppPorts.HTTP))
	}
	if methodFlag == rabbitMQFlag {
		connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", env.RabbitMQ.User, env.RabbitMQ.User, env.RabbitMQ.Host, env.RabbitMQ.Port)
		log.Default().Println("rabbitMQ connection string: ", connectionString)
		conn, err := amqp.Dial(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		log.Default().Println("rabbitMQ successfully connected: ", !conn.IsClosed())

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

		consumer := a.NewConsumer(q, c, rg, dpg)
		log.Printf("consumer running on queue %s", q.Name)
		consumer.Start(context.Background(), ch)
	} else {
		log.Fatalf("can't run application %s - please provide valid flag (app=setup|http|grpc|rabbitmq).", methodFlag)
	}
}
