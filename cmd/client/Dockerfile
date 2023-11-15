FROM golang:latest

WORKDIR /app

COPY . /app

RUN ls -la
RUN make deps

# TODO: adjust to make build correctly
RUN cp .env.example .env
RUN go build -o build/tcc-client-application cmd/client/main.go

# TODO: verify which por will run in Postman/Locust
EXPOSE 8082

CMD ["/app/build/tcc-client-application"]