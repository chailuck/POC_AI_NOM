// Makefile
.PHONY: build test run docker-build docker-run migrate-up migrate-down

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

run:
	go run cmd/server/main.go

docker-build:
	docker build -t tmf632-service .

docker-run:
	docker-compose up

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/tmf632db?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/tmf632db?sslmode=disable" down

lint:
	golangci-lint run

clean:
	rm -rf bin/
	go clean -cache


// Makefile additions
.PHONY: dev test-all generate-docs

dev:
	air

test-all:
	powershell -File scripts/test_all.ps1

generate-docs:
	swag init -g cmd/server/main.go -o api/swagger

lint:
	golangci-lint run