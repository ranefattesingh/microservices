FROM golang:1.26.6-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o user main.go

FROM alpine:latest

WORKDIR /app

EXPOSE 8080

COPY --from=0 /app/user /user

CMD ["/user"]