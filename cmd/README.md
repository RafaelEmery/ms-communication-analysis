## :bookmark_tabs: Main applications

The directory `cmd` contains the `main.go` and the `Dockerfile` for each application.

### *client*

Contains the main file for BFF. The `Dockerfile` is used to build on `docker-compose.yml`.

### *server*

Contains the main file that implements all server applications and it creates a new HTTP and gRPC servers and RabbitMQ queue connection. It's used a flag system to run specific applications:

```bash
# type=http|grpc|rabbitmq
go run -type=rabbitmq
```

Due to the approach of using flags instead of different `main.go` files for each application, the `Dockerfile` needed to be defined differently by using the sufix *.http*, *.grpc* and *rabbitmq* to be used to build on `docker-compose.yml`

:bell: It also has the `-type=setup` flag to run the Setup Application, a simple API to provide facilities to handle database during local tests, although can't be used inside a Docker container.

### *logprocesser*

Simple application to deal with logs from Docker containers. 