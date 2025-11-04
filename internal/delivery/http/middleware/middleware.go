package middleware

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"2025_2_404/internal/modules"
	"log"
	"net/http"
	"strings"
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

func (u *Middleware) Auth(next http.Handler) http.Handler {
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

func (u *Middleware) Peflite(next http.Handler) http.Handler {
    allowedOrigins := map[string]bool{
        "http://localhost:8000":        true,
        "http://127.0.0.1:8000":        true,
        "http://89.208.230.119:8000":   true,
    }

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
		log.Printf("Получен Origin: %s", origin)
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }

        // Эти заголовки можно устанавливать всегда — они не зависят от origin
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")
        w.Header().Set("Access-Control-Max-Age", "86400")

        // Если это предварительный запрос — сразу завершаем
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

