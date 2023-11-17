# Usage: make server-run type=setup|http|grpc|rabbitmq
server-run: 
	@echo "Running..."
	@go run cmd/server/main.go -type="${type}"

client-run: 
	@echo "Running..."
	@go run cmd/client/main.go

# @cp .env.example build/.env
server-build:
	@cp .env.example build/.env
	@go build -o build/tcc-server-application cmd/server/main.go

# @cp .env.example build/.env
client-build:
	@cp .env.example build/.env
	@go build -o build/tcc-client-application cmd/client/main.go

deps:
	go mod download
	go mod tidy

start:
	@docker-compose up -d

stop:
	@docker-compose down

build-and-start:
	@docker-compose up --build

restart: 
	@docker-compose restart

proto-generate:
	@rm internal/apps/message_grpc.pb.go
	@rm internal/apps/message.pb.go
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/apps/message.proto

remove-temporary-files:
	@rm -rf .tmp/*

load-testing:
	locust --host=http://localhost:3002
