FROM golang:tip-alpine3.22

WORKDIR /app

COPY . .

RUN go mod download

RUN go build main.go

EXPOSE 8080

CMD ["./main"]