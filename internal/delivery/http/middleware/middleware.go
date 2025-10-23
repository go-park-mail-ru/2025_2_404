package middleware

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"net/http"
	"strings"
	"2025_2_404/internal/modules"
)

type tokenUsecaseI interface {
	ValidateToken(tokenString string) (modeluser.ID, error)
}

type Middleware struct {
	tokenUsecase	tokenUsecaseI
}

func New(tokenUsecase tokenUsecaseI) *Middleware{
	return &Middleware{
		tokenUsecase: tokenUsecase,
	}
}

const (
	AuthHeaderKey  = "Authorization"
    AuthTypeBearer = "Bearer"
)

func (u *Middleware) AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthHeaderKey)
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != AuthTypeBearer {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		userID, err := u.tokenUsecase.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := modules.Set(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}