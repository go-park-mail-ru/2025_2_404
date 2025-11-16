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

FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

# RUN go build main.go

RUN go build -o /app/main ./cmd/ad

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/.env .

COPY --from=builder /app/keys ./keys

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["./main"]