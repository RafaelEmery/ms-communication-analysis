package main

import (
	"fmt"
	"log"

	goenv "github.com/Netflix/go-env"
	"github.com/RafaelEmery/performance-analysis-server/internal/apps/client"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // Load .env file
)

type Env struct {
	ClientPort string `env:"BFF_APP_PORT"`
	HTTPPort   string `env:"HTTP_APP_PORT"`
	GRPCPort   string `env:"GRPC_APP_PORT"`
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
	bff := client.NewBFFApp(fmt.Sprintf("http://localhost:%s", env.HTTPPort), fmt.Sprintf("localhost:%s", env.GRPCPort))
	bff.Routes(app)

	log.Default().Println("client application working")
	app.Listen(fmt.Sprintf(":%s", env.ClientPort))
}
