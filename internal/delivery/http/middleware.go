package handlers

import (
	"2025_2_404/pkg"
	"context"
	"crypto/ecdsa"
	"net/http"
	"strings"
)

type key string

const (
	UserIDKey key = "userID"
)

func AuthMiddleware(jwtPublicKey *ecdsa.PublicKey, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		userID, err := pkg.ValidateToken(jwtPublicKey, tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}