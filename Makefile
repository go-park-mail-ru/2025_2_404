generate-proto-ad:
	protoc -I protos/proto protos/proto/ad/ad.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative

generate-proto-storage:
	protoc -I protos/proto protos/proto/storage/storage.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative

generate-proto-slot:
	protoc -I protos/proto \
  protos/proto/slot/slot.proto \
  --go_out=./protos/gen/go/ \
  --go_opt=paths=source_relative \
  --go-grpc_out=./protos/gen/go/ \
  --go-grpc_opt=paths=source_relative

generate-proto-profile:
	protoc -I protos/proto \
  protos/proto/profile/profile.proto \
  --go_out=./protos/gen/go/ \
  --go_opt=paths=source_relative \
  --go-grpc_out=./protos/gen/go/ \
  --go-grpc_opt=paths=source_relative

generate-proto-auth:
	protoc -I protos/proto \
  protos/proto/auth/auth.proto \
  --go_out=./protos/gen/go/ \
  --go_opt=paths=source_relative \
  --go-grpc_out=./protos/gen/go/ \
  --go-grpc_opt=paths=source_relative

up-service-ad:
	go run cmd/ad/main.go

db-docker:
	docker compose up db --build -d

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...