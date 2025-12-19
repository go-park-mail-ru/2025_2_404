# FROM golang:alpine AS builder

# WORKDIR /app

# COPY . .

# RUN go mod download

# # RUN go build main.go

# RUN go build -o /app/main ./cmd/storage

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/.env .

# COPY internal/service/auth/keys ./internal/service/auth/keys

# COPY --from=builder /app/main ./main

# EXPOSE 8078

# CMD ["./storage"]

# FROM golang:1.24-alpine AS builder

# WORKDIR /build

# COPY ../../../go.mod ../../../go.sum ./

# COPY internal/delivery/grpc/storage ./internal/delivery/grpc/storage
# COPY internal/delivery/grpc/interceptor ./internal/delivery/grpc/interceptor
# COPY internal/service/storage ./internal/service/storage
# COPY protos ./protos
# COPY cmd/storage ./cmd/storage
# RUN CGO_ENABLED=0 GOOS=linux go build -o /storage ./cmd/storage/main.go

# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /root/

# COPY .env /root/.env
# COPY --from=builder /storage .
# RUN chmod +x ./storage
# ENV ENV_FILE=/root/.env

# EXPOSE 8078
# CMD ["./storage"]
# FROM  golang:1.24-alpine AS builder

# WORKDIR /app


# COPY go.mod go.sum ./
# RUN go mod download


# COPY . .


# RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/storage ./cmd/storage


# FROM --platform=linux/amd64 alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /app

# COPY --from=builder /app/bin/storage ./storage
# RUN chmod +x ./storage

# EXPOSE 8078

# ENTRYPOINT ["./storage"]

# 1. BUILDER STAGE
# Добавляем --platform=linux/amd64, чтобы сборка шла в "стандартном" режиме
FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o storage_app ./cmd/storage

# 2. RUNNER STAGE
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

# Копируем .env
COPY .env /root/.env

# Копируем бинарник
COPY --from=builder /build/storage_app .

# Копируем папку storage в корень → будет /storage/ad/usualads.jpg
COPY --from=builder /build/storage /storage/

# Устанавливаем переменную окружения
ENV ENV_FILE=/root/.env

# Даём права (на всякий случай)
RUN chmod +x ./storage_app

EXPOSE 8078

CMD ["./storage_app"]