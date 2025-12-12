    FROM golang:1.24-alpine AS builder

    WORKDIR /app

    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .

    RUN go build -o gateway-app ./cmd/main.go

    FROM alpine:latest
    
    WORKDIR /root/
    
    # (Опционально) Если нужны SSL сертификаты для внешних запросов
    # RUN apk --no-cache add ca-certificates

    COPY template ./template
    COPY .env ./.env
    COPY --from=builder /app/gateway-app .
    
    # Если у вас есть .env файл или папка configs, их нужно скопировать тоже
    # COPY --from=builder /app/.env . 
    # COPY --from=builder /app/configs ./configs
    
    EXPOSE 8080

    CMD ["./gateway-app"]