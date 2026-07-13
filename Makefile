EXECUTABLE := ecommerceplatform

.PHONY: build start clean

build:
	go build -o bin/${EXECUTABLE} cmd/api/main.go

start: build
	./bin/${EXECUTABLE}

clean:
	rm -rf bin/
