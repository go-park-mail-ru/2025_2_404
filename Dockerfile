FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

# RUN go build main.go

RUN go build -o /app/main ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/.env .

COPY --from=builder /app/keys ./keys

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["./main"]