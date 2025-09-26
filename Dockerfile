FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build main.go

FROM alpine:latest

WORKDIR /app

COPY .env .

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["./main"]