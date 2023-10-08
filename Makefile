run: 
	@echo "Running..."
	@go run cmd/main.go

deps:
	go mod download
	go mod tidy

start:
	@docker-compose up -d

stop:
	@docker-compose down

restart: 
	@docker-compose restart