package pkg

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// payload
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
} 

func GenerateToken(UserID, secretKey string) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject: "somebody",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	ss, err:= token.SignedString([]byte(secretKey))
	return ss, err
}

