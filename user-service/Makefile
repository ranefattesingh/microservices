migrate-up:
	 migrate -path ./repository/postgres/migrations -database "postgres://postgres:postgres@localhost:5432/users?sslmode=disable" up
create:
	migrate create -ext sql -dir repository/postgres/migrations -seq $(name)

sqlc:
	sqlc generate --file ./repository/postgres/sqlc/sqlc.yml

start:
	go run main.go