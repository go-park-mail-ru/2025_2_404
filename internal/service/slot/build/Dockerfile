FROM golang:1.24-alpine AS builder

WORKDIR /build

# Копируем зависимости
COPY go.mod go.sum ./

# Копируем только нужные пакеты (как в ad-service)
COPY . .

# Собираем бинарник
RUN go build -o /slot ./cmd/slot/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем .env (если используется) — но лучше через environment
# COPY .env /root/.env  # ← НЕ рекомендуется в продакшене
COPY .env /root/.env
COPY --from=builder /slot .
ENV ENV_FILE=/root/.env

# EXPOSE порт (для документации)
EXPOSE 8081

CMD ["./slot"]