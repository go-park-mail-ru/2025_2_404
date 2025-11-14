package pkg

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"crypto/ecdsa"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type UseCase struct {
	privateKey	*ecdsa.PrivateKey
	publicKey	*ecdsa.PublicKey
}

func New(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) *UseCase {
	return &UseCase{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

type Claims struct {
	UserID modeluser.ID `json:"user_id"`
	jwt.RegisteredClaims
} 

func (u *UseCase) GenerateToken(userID modeluser.ID) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err:= token.SignedString(u.privateKey)
	return ss, err
}

func (u *UseCase) ValidateToken(tokenString string) (modeluser.ID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return u.publicKey, nil
	})

	if err != nil {
		return modeluser.ID{}, fmt.Errorf("token parsing error: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return modeluser.ID{}, fmt.Errorf("invalid token")
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