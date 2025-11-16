generate-proto-ad:
	protoc -I protos/proto protos/proto/ad/ad.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative

up-service-ad:
	go run cmd/ad/main.go

db-docker:
	docker compose up db --build -d