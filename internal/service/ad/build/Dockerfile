# FROM golang:alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# WORKDIR /app/cmd/ad

# RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /root/

# COPY --from=builder /app/.env .

# COPY --from=builder /app/cmd/ad/main .

# EXPOSE 8081

# CMD ["./main"]

# FROM golang:alpine AS builder

# WORKDIR /app

# COPY . .

# RUN go mod download

# # RUN go build main.go

# RUN go build -o /app/main ./cmd/ad

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/.env .

# COPY internal/service/auth/keys ./internal/service/auth/keys

# COPY --from=builder /app/main ./main

# EXPOSE 8079

# CMD ["./ad"]

FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY ../../../go.mod ../../../go.sum ./

COPY internal/delivery/grpc/ad ./internal/delivery/grpc/ad
COPY internal/delivery/grpc/interceptor ./internal/delivery/grpc/interceptor
COPY internal/service/ad ./internal/service/ad
COPY protos ./protos
COPY cmd/ad ./cmd/ad    
RUN go build -o /ad ./cmd/ad/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY .env /root/.env
COPY --from=builder /ad .
ENV ENV_FILE=/root/.env

EXPOSE 8079
CMD ["./ad"]