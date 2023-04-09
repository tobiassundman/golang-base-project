.DEFAULT_GOAL := help

PROJECT_NAME := go-demo-app

.PHONY: help
help:
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tools
make tools: ## Install required tools
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/go-bindata/go-bindata/go-bindata@latest

.PHONY: build
build: check test build_image ## Checks, tests and builds the docker image

.PHONY: format
format: ## Format code
	gofmt -s -w .

.PHONY: check
check: ## Runs code checks
	go vet ./...
	staticcheck ./...

.PHONY: vulncheck
vulncheck: ## Runs code checks
	govulncheck ./...

.PHONY: test
test: ## Runs unit tests
	go test ./...

.PHONY: build_image
build_image: ## Builds docker image
	docker image build . --file build/demo-app/Dockerfile -t go-demo-app:latest

.PHONY: package-migrations
package-migrations: ## Packages migrations into a go file
	go-bindata -pkg migrations -ignore migrations.go -nometadata -prefix db/migrations/ -o db/migrations/migrations.go ./db/migrations/

.PHONY: deploy_up
deploy_up: ## Starts the application in docker-compose
	docker-compose -f deployments/docker-compose.yml up -d 

.PHONY: deploy_down
deploy_down: ## Stops the application in docker-compose
	docker-compose -f deployments/docker-compose.yml down 

.PHONY: migrate
migrate: ## Stops the application in docker-compose
	go run cmd/db-migration/main.go