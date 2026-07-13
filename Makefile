EXECUTABLE := ecommerceplatform

DB_URL := postgres://postgres:postgres@10.255.255.254:5432/ecommerce?sslmode=disable
MIGRATIONS_DIR := migrations

.PHONY: build start clean create-migration migrate-up migrate-down migrate-force migrate-version


build:
	go build -o bin/${EXECUTABLE} cmd/api/main.go

start: build
	./bin/${EXECUTABLE}

clean:
	rm -rf bin/

create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create-migration name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down 1

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Usage: make migrate-force version=3"; \
		exit 1; \
	fi
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) force $(version)

migrate-version:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) version
