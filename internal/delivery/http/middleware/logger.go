package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// responseWriterWithStatus — перехватывает статус и размер ответа
type responseWriterWithStatus struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func (rw *responseWriterWithStatus) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWithStatus) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	if !rw.wroteHeader {
		rw.status = http.StatusOK
		rw.wroteHeader = true
	}
	rw.size += size
	return size, err
}

func (rw *responseWriterWithStatus) Status() int {
	if !rw.wroteHeader {
		return http.StatusOK
	}
	return rw.status
}

func (rw *responseWriterWithStatus) Size() int {
	return rw.size
}

// AccessLogMiddleware — middleware для логирования с log/slog
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()[:8]
		}

		// IP с поддержкой прокси
		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}

		// Оборачиваем ResponseWriter
		rw := &responseWriterWithStatus{
			ResponseWriter: w,
			status:         http.StatusOK,
			size:           0,
			wroteHeader:    false,
		}

		// Добавляем request_id в контекст (если понадобится в handler'ах)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		newReq := r.WithContext(ctx)

		// Выполняем обработчик
		next.ServeHTTP(rw, newReq)

		duration := time.Since(start)
		durationMs := float64(duration.Microseconds()) / 1000.0

		// Базовые атрибуты
		attrs := []slog.Attr{
			slog.String("req_id", requestID),
			slog.String("method", r.Method),
			slog.String("url", r.URL.EscapedPath()),
			slog.String("query", r.URL.RawQuery),
			slog.Int("status", rw.Status()),
			slog.Int("size", rw.Size()),
			slog.String("ip", remoteAddr),
			slog.String("ua", r.UserAgent()),
			slog.Float64("took_ms", durationMs),
		}

		// Выбираем уровень и сообщение
		var level slog.Level
		var msg string
		switch {
		case rw.Status() >= 500:
			level = slog.LevelError
			msg = "⚠️ SERVER ERROR"
		case rw.Status() >= 400:
			level = slog.LevelWarn
			msg = "⚠️ CLIENT ERROR"
		default:
			level = slog.LevelInfo
			msg = "✅ OK"
		}

		// Логируем
		slog.LogAttrs(r.Context(), level, msg, attrs...)
	})
}
