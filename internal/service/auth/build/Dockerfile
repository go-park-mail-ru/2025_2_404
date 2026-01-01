FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /auth ./cmd/auth/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /auth .
COPY internal/service/auth/keys ./internal/service/auth/keys

EXPOSE 8077
CMD ["./auth"]