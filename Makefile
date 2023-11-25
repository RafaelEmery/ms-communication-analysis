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

# Usage: make start service=service-name (optional)
start:
	@docker-compose up -d ${service}

# fix-start:
# 	@docker-compose up -d bff-app

stop:
	@docker-compose down

get-container-logs:
	@docker-compose logs bff-app
	@docker-compose logs ${service}-app

# Usage: make start service=service-name (optional)
start-with-build:
	@docker-compose up -d --build ${service}

docker-client-build:
	docker build -t tcc-client-application -f cmd/client/Dockerfile .

proto-generate:
	@rm internal/apps/message_grpc.pb.go
	@rm internal/apps/message.pb.go
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/apps/message.proto

remove-temporary-files:
	@rm -rf .tmp/*

load-testing:
	locust --host=http://0.0.0.0:3002

# Usage: make get-logs service=server-service-name (without "-app") index=index
get-logs:
	@docker-compose logs ${service}-app > logs/${service}-app-${index}.txt

remove-containers-logs:
	@rm -rf logs/*

# Usage: make proccess-logs service=server-service-name (without "-app") index=index
proccess-logs:
	go run cmd/logprocesser/main.go ./logs/${service}-app-${index}.txt

start-specific:
	./scripts/start-specific.sh ${service}

restart-script:
	./scripts/restart-all.sh

specific-restart-script:
	./scripts/restart-specific.sh ${service}

extract-script:
	./scripts/extract.sh ${service} ${index}