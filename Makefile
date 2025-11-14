proto-auth:
	protoc -I proto protos/auth/auth.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/auth/auth.proto --go-grpc_opt=paths=source_relative

	protoc \
  -I protos \
  protos/auth/auth.proto \
  --go_out=protos \
  --go_opt=paths=source_relative \
  --go-grpc_out=protos \
  --go-grpc_opt=paths=source_relative
