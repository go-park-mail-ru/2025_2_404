package middleware

import (
	"net/http"
	"strconv"
	"time"

	"2025_2_404/pkg/metrics"
)

// MetricsMiddleware — middleware для сбора Prometheus метрик.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Увеличиваем счётчик активных запросов
		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		start := time.Now()

		// Нормализуем path для метрик (убираем динамические ID)
		path := normalizePath(r.URL.Path)

		// Оборачиваем ResponseWriter для перехвата статуса и размера
		rw := &metricsResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
			size:           0,
		}

		// Выполняем обработчик
		next.ServeHTTP(rw, r)

		// Записываем метрики
		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(rw.status)

		// Счётчик запросов
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, statusStr).Inc()

		// Гистограмма времени ответа
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		// Гистограмма размера ответа
		metrics.HTTPResponseSize.WithLabelValues(r.Method, path).Observe(float64(rw.size))

		// Счётчик по классам статусов
		statusClass := metrics.GetStatusClass(rw.status)
		metrics.HTTPStatusCodes.WithLabelValues(statusClass).Inc()
	})
}

// metricsResponseWriter перехватывает статус и размер ответа.
type metricsResponseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.status = http.StatusOK
		rw.wroteHeader = true
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// normalizePath нормализует URL path, заменяя UUID и числовые ID на плейсхолдеры.
// Это предотвращает "взрыв" кардинальности метрик.
func normalizePath(path string) string {
	// Простая нормализация: заменяем UUID и числа в пути
	// /api/ads/123 -> /api/ads/:id
	// /api/slots/abc-def-123 -> /api/slots/:id

	segments := splitPath(path)
	for i, seg := range segments {
		if isUUID(seg) || isNumeric(seg) {
			segments[i] = ":id"
		}
	}
	return joinPath(segments)
}

func splitPath(path string) []string {
	var segments []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				segments = append(segments, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		segments = append(segments, path[start:])
	}
	return segments
}

func joinPath(segments []string) string {
	if len(segments) == 0 {
		return "/"
	}
	result := ""
	for _, seg := range segments {
		result += "/" + seg
	}
	return result
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isUUID(s string) bool {
	// UUID формат: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 символов)
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
