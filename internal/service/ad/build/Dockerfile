

FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY ../../../go.mod ../../../go.sum ./

COPY . .

RUN go build -o /ad ./cmd/ad/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /ad .

EXPOSE 8079
CMD ["./ad"]