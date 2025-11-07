.PHONY: help build run docker-build docker-run clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the Go application
	go build -o go-ghant

run: ## Run the application locally
	go run .

docker-build: ## Build Docker container
	docker build -t go-ghant .

docker-run: ## Run Docker container
	docker run -p 8080:8080 -v $(CURDIR)/data:/root go-ghant

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose
	docker-compose down

clean: ## Clean build artifacts
	rm -f go-ghant charts.json
	rm -rf data/

test: ## Run tests
	go test -v ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

example: ## Copy example data
	cp charts.example.json charts.json
