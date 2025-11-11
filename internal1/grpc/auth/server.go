package auth

import (
	adnetv1 "2025_2_404/protos/gen/go/adnet"
	"context"

	"google.golang.org/grpc"
) 

type serverAPI struct {
	adnetv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	adnetv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *adnetv1.LoginRequest) (*adnetv1.LoginResponse, error) {
	panic("init failed")
}