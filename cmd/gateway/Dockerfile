    FROM golang:1.24-alpine AS builder

    WORKDIR /app

    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .

    RUN go build -o gateway-app ./cmd/main.go

    FROM alpine:latest
    
    WORKDIR /root/
    

    COPY template ./template
    COPY .env ./.env
    COPY --from=builder /app/gateway-app .
    
    EXPOSE 8080

    CMD ["./gateway-app"]