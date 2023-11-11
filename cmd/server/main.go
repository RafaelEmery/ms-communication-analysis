package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	goenv "github.com/Netflix/go-env"
	"google.golang.org/grpc"

	i "github.com/RafaelEmery/performance-analysis-server/internal"
	a "github.com/RafaelEmery/performance-analysis-server/internal/apps/server"
	u "github.com/RafaelEmery/performance-analysis-server/internal/usecases"
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // Load .env file
	_ "github.com/lib/pq"                 // Load postgres connection
)

const (
	setupFlag = "setup"
	httpFlag  = "http"
	grpcFlag  = "grpc"
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

	c := u.NewCreateUseCase(r)
	rg := u.NewReportUseCase(r)
	dpg := u.NewGetByDiscountUseCase(r)

	flag.Parse()
	log.Default().Println("running app by flag", flag.Arg(0))

	if flag.Arg(0) == setupFlag {
		app := fiber.New()
		setupApp := a.NewSetupApp(context.Background(), r)
		setupApp.Routes(app)

		log.Default().Println("setup application working")
		app.Listen(fmt.Sprintf(":%s", env.AppPorts.Setup))
	}
	if flag.Arg(0) == grpcFlag {
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
	if flag.Arg(0) == httpFlag {
		app := fiber.New()
		httpApp := a.NewHttpApp(context.Background(), c, rg, dpg)
		httpApp.Routes(app)

		log.Default().Println("HTTP endpoints working")
		app.Listen(fmt.Sprintf(":%s", env.AppPorts.HTTP))
	} else {
		log.Fatalf("can't run application %s - please provide valid flag (app=setup|http|grpc).", flag.Arg(0))
	}
}
