# Usage: make server-run app=setup|http|grpc|rabbitmq
server-run: 
	@echo "Running..."
	@go run cmd/server/main.go ${app}

client-run: 
	@echo "Running..."
	@go run cmd/client/main.go

# TODO: fix make build - make: 'build' está atualizado.
server-build:
	@go build -o build/client cmd/client/main.go

client-build:
	@go build -o build/server cmd/server/main.go

deps:
	go mod download
	go mod tidy

start:
	@docker-compose up -d

docker-server-build:
	docker build -t server -f cmd/server/Dockerfile .

# Usage: make docker-server-run flag=setup|http|grpc|rabbitmq
# TODO: fix flag on server main.go and dockerfile
docker-server-run:
	docker run -p 8081:8081 server 

docker-client-build:
	docker build -t tcc-client-application -f cmd/client/Dockerfile .

docker-client-run:
	docker run -p 8082:8082 tcc-client-application

stop:
	@docker-compose down

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