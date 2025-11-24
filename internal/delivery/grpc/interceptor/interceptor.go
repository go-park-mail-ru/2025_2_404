package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
    

	"github.com/google/uuid"
	authProto "2025_2_404/protos/auth"
)

type ctxKey string
const UserIDKey ctxKey = "userID"

func AuthInterceptor(authClient authProto.AuthClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
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
			return nil, status.Error(codes.Unauthenticated, "invalid auth format")
		}

		resp, err := authClient.ValidateToken(ctx, &authProto.TokenRequest{
            Token: token,
        })
        
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token validation failed: %v", err)
		}

        userID, err := uuid.Parse(resp.GetUserId())
        if err != nil {
            return nil, status.Errorf(codes.Internal, "invalid user id from auth service")
        }

		newCtx := context.WithValue(ctx, UserIDKey, userID)
		return handler(newCtx, req)
	}
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
    val := ctx.Value(UserIDKey)
    if val == nil {
        return uuid.Nil, status.Error(codes.Unauthenticated, "user id not found in context")
    }
    
    id, ok := val.(uuid.UUID)
    if !ok {
        return uuid.Nil, status.Error(codes.Internal, "user id is of wrong type")
    }
    return id, nil
}