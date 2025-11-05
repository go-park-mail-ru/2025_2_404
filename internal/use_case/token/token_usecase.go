package token

import (
	"2025_2_404/internal/config"
	modeluser "2025_2_404/internal/domain/models/user"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UseCase struct {
	privateKey	*ecdsa.PrivateKey
	publicKey	*ecdsa.PublicKey
}

func New(cfg *config.Config) *UseCase {
	return &UseCase{
		privateKey: cfg.AppConfig.JwtPrivateKey,
		publicKey: cfg.AppConfig.JwtPublicKey,
	}
}

type Claims struct {
	UserID modeluser.ID `json:"user_id"`
	jwt.RegisteredClaims
} 

func (usecase *UseCase) GenerateToken(userID modeluser.ID) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject: "somebody",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	ss, err:= token.SignedString(usecase.privateKey)
	return ss, err
}

func (usecase *UseCase) ValidateToken(tokenString string) (modeluser.ID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func (token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return usecase.publicKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("token parsing error: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, fmt.Errorf("invalid token")
}

func (usecase *UseCase) InvalidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func (token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return "", fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return usecase.publicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("token parsing error: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now())
		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

		ss, err:= token.SignedString(usecase.privateKey)
		return ss, err
	}

	return "", fmt.Errorf("invalid token")
}