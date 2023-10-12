package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	goenv "github.com/Netflix/go-env"

	i "github.com/RafaelEmery/performance-analysis-server/internal"
	a "github.com/RafaelEmery/performance-analysis-server/internal/apps"
	u "github.com/RafaelEmery/performance-analysis-server/internal/usecases"
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // Load .env file
	_ "github.com/lib/pq"                 // Load postgres connection
)

type Env struct {
	AppPort string `env:"APP_PORT"`
	DB      struct {
		Driver   string `env:"DB_DRIVER"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
	}
}

func getEnv() (*Env, error) {
	env := &Env{}
	_, err := goenv.UnmarshalFromEnviron(env)

	log.Default().Println("environment loaded")

	return env, err
}

func getDatabase(env *Env) (*sql.DB, error) {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DB.Host, env.DB.Port, env.DB.User, env.DB.Password, env.DB.Name)

	db, err := sql.Open(env.DB.Driver, dbSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Default().Println("database connected")

	return db, nil
}

func setupHttpApp(app *fiber.App, c a.Creator, rg a.ReportGenerator) {
	httpApp := a.NewHttpApp(context.Background(), c, rg)
	httpApp.Routes(app)

	log.Default().Println("HTTP endpoints working")
}

func main() {
	env, err := getEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := getDatabase(env)
	if err != nil {
		log.Fatal(err)
	}

	r := i.NewRepository(db)
	app := fiber.New()

	creator := u.NewCreateUseCase(r)
	reportGenerator := u.NewReportUseCase(r)

	setupHttpApp(app, creator, reportGenerator)

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON("Testing...")
	})

	app.Listen(env.AppPort)
}
