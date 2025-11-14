package interceptor

import (
    "context"
    "strings"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
    modeluser "2025_2_404/internal/service/auth/domain"
)

type tokenI interface {
    ValidateToken(tokenStr string) (modeluser.ID, error)
}

func AuthInterceptor(tokenI tokenI) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        if strings.HasSuffix(info.FullMethod, "/Register") ||
           strings.HasSuffix(info.FullMethod, "/Login") {
            return handler(ctx, req)
        }

        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "no metadata")
        }

        auth := md.Get("authorization")
        if len(auth) == 0 {
            return nil, status.Error(codes.Unauthenticated, "no token")
        }

        token := strings.TrimPrefix(auth[0], "Bearer ")
        if token == auth[0] {
            return nil, status.Error(codes.Unauthenticated, "invalid format")
        }

        userID, err := tokenI.ValidateToken(token)
        if err != nil {
            return nil, status.Error(codes.Unauthenticated, "invalid token")
        }

        ctx = context.WithValue(ctx, "userID", userID)
        return handler(ctx, req)
    }
}