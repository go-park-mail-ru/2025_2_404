package auth

import (
	adnetv1 "2025_2_404/protos/gen/go/adnet"
	"google.golang.org/grpc"
) 

type serverAPi struct {
	adnetv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	
}