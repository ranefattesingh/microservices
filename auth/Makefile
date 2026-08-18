build:
	@go build -o bin/user main.go

start:
	./bin/user

migrate-up:
	migrate \
    -path repository/db/migrations \
    -database "postgres://test:test1234@localhost:5432/auth?sslmode=disable" \
    up