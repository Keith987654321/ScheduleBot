FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o application .

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/application .

ENV DB_HOST = \
    DB_USERNAME= \
    DB_PASSWORD= \ 
    DB_NAME= \
    SSL_MODE= \
    TOKEN=

CMD ["./application"]
