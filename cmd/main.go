package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	goenv "github.com/Netflix/go-env"
	"google.golang.org/grpc"

	i "github.com/RafaelEmery/performance-analysis-server/internal"
	a "github.com/RafaelEmery/performance-analysis-server/internal/apps"
	u "github.com/RafaelEmery/performance-analysis-server/internal/usecases"
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // Load .env file
	_ "github.com/lib/pq"                 // Load postgres connection
)

type Env struct {
	AppPorts struct {
		HTTP string `env:"HTTP_APP_PORT"`
		GRPC string `env:"GRPC_APP_PORT"`
	}
	DB struct {
		Driver   string `env:"DB_DRIVER"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_DATABASE"`
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
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

	db, err := connectDatabase(env)
	if err != nil {
		log.Fatal(err)
	}

	r := i.NewRepository(db)
	app := fiber.New()

	c := u.NewCreateUseCase(r)
	rg := u.NewReportUseCase(r)
	dpg := u.NewGetByDiscountUseCase(r)

	// TODO: use flags commands to specify which server is to run
	// TODO: separate main files to HTTP app and gRPC app
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

	setupApp := a.NewSetupApp(context.Background(), r)
	setupApp.Routes(app)
	log.Default().Println("setup application working")

	httpApp := a.NewHttpApp(context.Background(), c, rg, dpg)
	httpApp.Routes(app)
	log.Default().Println("HTTP endpoints working")

	app.Listen(fmt.Sprintf(":%s", env.AppPorts.HTTP))
}
