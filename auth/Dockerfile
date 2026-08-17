# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct
COPY . .
RUN go build -o /auth ./...

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /auth /usr/local/bin/auth
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/auth"]
