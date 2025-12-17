// Package service provides JWT token generation and validation services.
package service

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"2025_2_404/pkg/globalerrors"
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UseCaseJWT struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	logger     *zap.Logger
}

func NewJWT(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, logger *zap.Logger) *UseCaseJWT {
	return &UseCaseJWT{
		privateKey: privateKey,
		publicKey:  publicKey,
		logger:     logger,
	}
}

type Claims struct {
	UserID modeluser.ID `json:"user_id"`
	jwt.RegisteredClaims
}

func (u *UseCaseJWT) GenerateToken(userID modeluser.ID) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(u.privateKey)
	if err != nil {
		return "", globalerrors.ErrInternal
	}
	return ss, nil
}

func (u *UseCaseJWT) ValidateToken(ctx context.Context, tokenString string) (modeluser.ID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return u.publicKey, nil
	})

	if err != nil {
		return modeluser.ID{}, globalerrors.ErrNoAuth
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return modeluser.ID{}, globalerrors.ErrNoAuth
}

// func (u *UseCase) InvalidateToken(tokenString string) (string, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
// 			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
// 		}
// 		return u.publicKey, nil
// 	})
// 	if err != nil || !token.Valid{

// 	}
// }
