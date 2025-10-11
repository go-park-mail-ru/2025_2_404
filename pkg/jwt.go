package pkg

import (
	"2025_2_404/internal/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// payload
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
} 

var cfg = config.GetAppConfig()
var secretKey = cfg.JwtToken

func GenerateToken(UserID int64) (string, error) {
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

func ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func (token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("Ошибка парсинга токена")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, fmt.Errorf("невалидный токен")
}

