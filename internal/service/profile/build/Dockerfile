FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY ../../../go.mod ../../../go.sum ./

COPY internal/delivery/grpc/profile ./internal/delivery/grpc/profile
COPY internal/delivery/grpc/interceptor ./internal/delivery/grpc/interceptor
COPY internal/service/profile ./internal/service/profile
COPY protos ./protos
COPY cmd/profile ./cmd/profile

RUN go build -o /profile ./cmd/profile/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
RUN apk add --no-cache curl
WORKDIR /root/

COPY .env /root/.env
COPY --from=builder /profile .
ENV ENV_FILE=/root/.env

EXPOSE 8076
CMD ["./profile"]