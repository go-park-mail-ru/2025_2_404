FROM golang:alpine AS builder

WORKDIR cmd/app

COPY . .

RUN go mod download

# RUN go build main.go

RUN go build -o /app/main ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY .env .

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["./main"]