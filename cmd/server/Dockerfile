FROM golang:latest

WORKDIR /app

COPY . /app

RUN ls -la
RUN make deps

# TODO: adjust to make build correctly
RUN cp .env.example .env
RUN go build -o build/server cmd/server/main.go

# TODO: verify which por will run in Postman/Locust
EXPOSE 8081

# Passing flag to server command and setting http as default
CMD ["/app/build/server", "-type", "http"]