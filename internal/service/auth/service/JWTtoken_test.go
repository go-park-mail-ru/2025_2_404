package service

import (
	"2025_2_404/pkg/globalerrors"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func generateTestKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func TestGenerateToken_Success(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	userID := uuid.New()
	token, err := useCase.GenerateToken(userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestGenerateToken_TokenStructure(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	userID := uuid.New()
	tokenString, err := useCase.GenerateToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		t.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Error("expected token to be valid")
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, claims.UserID)
	}

	now := time.Now()
	if claims.IssuedAt.Time.After(now) {
		t.Error("IssuedAt should not be in the future")
	}

	expectedExpiry := now.Add(24 * time.Hour)
	if claims.ExpiresAt.Time.Before(now) {
		t.Error("ExpiresAt should be in the future")
	}

	if claims.ExpiresAt.Time.After(expectedExpiry.Add(1 * time.Minute)) {
		t.Error("ExpiresAt is too far in the future")
	}
}

func TestValidateToken_Success(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	userID := uuid.New()
	token, err := useCase.GenerateToken(userID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	validatedUserID, err := useCase.ValidateToken(context.Background(), token)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, validatedUserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	invalidToken := "invalid.token.string"

	_, err = useCase.ValidateToken(context.Background(), invalidToken)
	if err != globalerrors.ErrNoAuth {
		t.Errorf("expected ErrNoAuth, got %v", err)
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	userID := uuid.New()
	expiredTime := time.Now().Add(-1 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(expiredTime.Add(-24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	expiredTokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = useCase.ValidateToken(context.Background(), expiredTokenString)
	if err != globalerrors.ErrNoAuth {
		t.Errorf("expected ErrNoAuth for expired token, got %v", err)
	}
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hmacToken, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign HMAC token: %v", err)
	}

	_, err = useCase.ValidateToken(context.Background(), hmacToken)
	if err != globalerrors.ErrNoAuth {
		t.Errorf("expected ErrNoAuth for wrong signing method, got %v", err)
	}
}

func TestValidateToken_WrongPublicKey(t *testing.T) {
	privateKey1, _, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys 1: %v", err)
	}

	_, publicKey2, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys 2: %v", err)
	}

	logger, _ := zap.NewDevelopment()

	useCaseSign := NewJWT(privateKey1, &privateKey1.PublicKey, logger)
	useCaseValidate := NewJWT(privateKey1, publicKey2, logger)

	userID := uuid.New()
	token, err := useCaseSign.GenerateToken(userID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = useCaseValidate.ValidateToken(context.Background(), token)
	if err != globalerrors.ErrNoAuth {
		t.Errorf("expected ErrNoAuth for wrong public key, got %v", err)
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	_, err = useCase.ValidateToken(context.Background(), "")
	if err != globalerrors.ErrNoAuth {
		t.Errorf("expected ErrNoAuth for empty token, got %v", err)
	}
}

func TestNewJWT(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	useCase := NewJWT(privateKey, publicKey, logger)

	if useCase == nil {
		t.Error("expected non-nil useCase")
	}

	if useCase.privateKey != privateKey {
		t.Error("expected privateKey to match")
	}

	if useCase.publicKey != publicKey {
		t.Error("expected publicKey to match")
	}

	if useCase.logger != logger {
		t.Error("expected logger to match")
	}
}
