include .env
export

# ======= DOCKER COMPOSE =======
up:
	docker compose up -d

down:
	docker compose down

restart: down up

build:
	docker compose build

logs:
	docker compose logs -f app-pr-reviewer-service

# ======= MIGRATIONS =======
migrate-up:
	migrate -path ./migrations/sql -database "postgres://$(STORAGES_POSTGRES_USER):$(STORAGES_POSTGRES_PASS)@$(STORAGES_POSTGRES_HOST):$(STORAGES_POSTGRES_PORT)/$(STORAGES_POSTGRES_NAME)?sslmode=$(STORAGES_POSTGRES_SSLM)" up

migrate-down:
	migrate -path ./migrations/sql -database "postgres://$(STORAGES_POSTGRES_USER):$(STORAGES_POSTGRES_PASS)@$(STORAGES_POSTGRES_HOST):$(STORAGES_POSTGRES_PORT)/$(STORAGES_POSTGRES_NAME)?sslmode=$(STORAGES_POSTGRES_SSLM)" down

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Error: укажи имя миграции через 'name=...'" && exit 1; \
	fi
	migrate create -ext sql -dir ./migrations/sql $(name)

# ======= LINT =======
lint:
	golangci-lint run ./cmd/... ./internal/... ./pkg/...

# lint-fix:
# 	golangci-lint run --fix ./cmd/... ./internal/... ./pkg/...

# ======= TESTS =======
test:
	go test -v -short ./internal/... ./pkg/...

test-coverage:
	mkdir -p ./tmp
	go test \
		-short \
		-count=1 \
		-race \
		-coverprofile=./tmp/coverage.out \
		./internal/... \
		./pkg/...
	@echo ""
	@echo "Coverage (excluding mocks):"
	@grep -v "/mocks" ./tmp/coverage.out | grep -v "mocks--exported" > ./tmp/coverage_no_mocks.out
	@go tool cover -func=./tmp/coverage_no_mocks.out | grep "total:" | awk '{print "Total coverage:", $$3}'
	go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.html


test-e2e:
	go test -v ./tests/e2e/...

# ======= MOCKS =======
gen-mocks:
	go generate ./internal/...

# ======= DEV =======
dev-run:
	go run ./cmd/pr-reviewer-service/main.go