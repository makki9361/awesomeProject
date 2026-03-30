.PHONY: run build rebuild docker-build docker-rebuild

BINARY_NAME=app
MIGRATIONS_PATH=./migrations
DB_URL=postgres://postgres:password@localhost:5432/myapp?sslmode=disable

run:
	go run cmd/app/main.go

build:
	go build -o bin/app cmd/app/main.go

rebuild: build
	./bin/app

docker-build:
	docker-compose build

docker-rebuild:
	docker-compose build --no-cache

docker-restart:
	docker-compose down
	docker-compose build
	docker-compose up -d

docker-reset:
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up -d	