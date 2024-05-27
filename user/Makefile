APP=user
APP_EXECUTABLE=./bin/$(APP)
DB_USER=postgres
COVER_FILE=coverage.out
COVER_HTML=coverage.html
UNIT_TEST_PKG_LIST=`go list ./... | grep -v vendor | grep -v integration-test`
GOBIN=$$(go env GOPATH)/bin

.PHONY: swagger

install.tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.58.2
	$(GOBIN)/golangci-lint --version
	go install github.com/swaggo/swag/cmd/swag@latest
	$(GOBIN)/swag --version
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GOBIN)/migrate --version

build: clean
	go build -o $(APP_EXECUTABLE) main.go

clean:
	rm -rf bin coverage.out coverage.html golangci-lint-report.txt

swagger:
	$(GOBIN)/swag init -g handler/user.go -o ./swagger -ot go

lint:
	$(GOBIN)/golangci-lint run -c .golangci.yml

lint.fix:
	$(GOBIN)/golangci-lint run -c .golangci.yml --fix

start:
	$(APP_EXECUTABLE)

test.unit:
	go test $(UNIT_TEST_PKG_LIST) -timeout 1h

test.coverage: clean
	ENVIRONMENT=test go test $(UNIT_TEST_PKG_LIST) -p=2 -v
	@go list ./... | grep -v "vendor" | grep -v "integration" | xargs go test -count 1 -cover -short -race -timeout 1m -coverprofile ${COVER_FILE}
	@go tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
	@go tool cover -func ${COVER_FILE} | tail -1 | xargs echo test coverage:

test.integration:
	go test `go list ./... | grep integration-test` -timeout 1m

test.benchmark:
	go test -v -bench=. `go list ./... | grep -v integration-test` -benchtime=10s -run notest -benchmem -timeout 1m

db.setup: db.create db.migrate

db.create:
	createdb -h localhost -U $(DB_USER) -O $(DB_USER) -Eutf8 $(DB_NAME)

db.migrate: $(APP_EXECUTABLE)
	$(APP_EXECUTABLE) migrate

db.rollback: $(APP_EXECUTABLE)
	$(APP_EXECUTABLE) rollback

db.drop:
	dropdb --if-exists -h localhost -U$(DB_USER) $(DB_NAME) 

db.reset: db.drop db.create db.migrate

db.create-migration:
	migrate create -dir ./repo/psql/migrations -ext sql $(MIGRATION_NAME)