FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY ../../../go.mod ../../../go.sum ./

COPY internal/delivery/grpc/auth ./internal/delivery/grpc/auth
COPY internal/delivery/grpc/interceptor ./internal/delivery/grpc/interceptor
COPY internal/service/auth ./internal/service/auth
COPY internal/pkg/jwt ./internal/pkg/jwt
COPY protos ./protos
COPY cmd/auth ./cmd/auth

RUN go build -o /auth ./cmd/auth/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY .env /root/.env
COPY --from=builder /auth .
COPY internal/service/auth/keys ./internal/service/auth/keys
ENV ENV_FILE=/root/.env

EXPOSE 50001
CMD ["./auth"]