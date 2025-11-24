generate-proto-ad:
	protoc -I protos/proto protos/proto/ad/ad.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative

generate-proto-storage:
	protoc -I protos/proto protos/proto/storage/storage.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative

up-service-ad:
	go run cmd/ad/main.go

db-docker:
	docker compose up db --build -d
proto-auth:
	protoc -I proto protos/auth/auth.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/auth/auth.proto --go-grpc_opt=paths=source_relative

	protoc \
  -I protos \
  protos/auth/auth.proto \
  --go_out=protos \
  --go_opt=paths=source_relative \
  --go-grpc_out=protos \
  --go-grpc_opt=paths=source_relative

	protoc \
  -I protos \
  protos/profile/profile.proto \
  --go_out=protos \
  --go_opt=paths=source_relative \
  --go-grpc_out=protos \
  --go-grpc_opt=paths=source_relative
  
