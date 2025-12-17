// Package interceptor provides gRPC interceptors for authentication and authorization.
package interceptor

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authProto "2025_2_404/protos/gen/go/auth"

	"github.com/google/uuid"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

var publicMethods = map[string]bool{
	"/slot.SlotServ/GetSlot":               true,
	"/slot.SlotServ/CreateMetric":          true,
	"/ad.AdServ/GetAdDetailForSlot":        true,
	"/ad.AdServ/GetAdSlot":                 true,
	"/profile.Profile/UpdatePaymentStatus": true,
	"/profile.Profile/AddBalance":          true,
}

func AuthInterceptor(authClient authProto.AuthClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if publicMethods[info.FullMethod] {
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
			return nil, status.Error(codes.Unauthenticated, "invalid auth format")
		}

		resp, err := authClient.ValidateToken(ctx, &authProto.TokenRequest{
			Token: token,
		})

		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token validation failed: %v", err)
		}
		fmt.Printf("DEBUG INTERCEPTOR: Auth returned UserID string: '%s'\n", resp.GetUserId())
		userID, err := uuid.Parse(resp.GetUserId())
		fmt.Printf("DEBUG INTERCEPTOR: Auth returned UserID string: '%s'\n", userID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "invalid user id from auth service")
		}

		newCtx := context.WithValue(ctx, UserIDKey, userID)
		fmt.Printf("DEBUG INTERCEPTOR: Auth returned UserID string: '%s'\n", newCtx)

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

func InitAuthInterceptor() (grpc.UnaryServerInterceptor, *grpc.ClientConn) {
	authAddr := os.Getenv("AUTH_ADDR")
	if authAddr == "" {
		authAddr = "auth_service:8077"
	}

	conn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("CRITICAL: Failed to connect to Auth Service: %v", err)
	}

	authClient := authProto.NewAuthClient(conn)
	return AuthInterceptor(authClient), conn
}
