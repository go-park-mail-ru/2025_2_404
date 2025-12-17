// Package pkg provides utility functions for common tasks such as sending JSON responses.
package pkg

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"data":    data,
	}); err != nil {
		slog.Warn("Failed to encode JSON response", "error", err)
	}
}
