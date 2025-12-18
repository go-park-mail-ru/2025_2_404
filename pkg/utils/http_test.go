package utils

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestHTTPStatusFromCode(t *testing.T) {
	tests := []struct {
		name         string
		code         codes.Code
		expectedHTTP int
	}{
		{"OK", codes.OK, http.StatusOK},
		{"Canceled", codes.Canceled, http.StatusRequestTimeout},
		{"Unknown", codes.Unknown, http.StatusBadRequest},
		{"InvalidArgument", codes.InvalidArgument, http.StatusBadRequest},
		{"DeadlineExceeded", codes.DeadlineExceeded, http.StatusRequestTimeout},
		{"NotFound", codes.NotFound, http.StatusNotFound},
		{"AlreadyExists", codes.AlreadyExists, http.StatusConflict},
		{"PermissionDenied", codes.PermissionDenied, http.StatusForbidden},
		{"Unauthenticated", codes.Unauthenticated, http.StatusUnauthorized},
		{"ResourceExhausted", codes.ResourceExhausted, http.StatusTooManyRequests},
		{"FailedPrecondition", codes.FailedPrecondition, http.StatusPreconditionFailed},
		{"Aborted", codes.Aborted, http.StatusConflict},
		{"OutOfRange", codes.OutOfRange, http.StatusBadRequest},
		{"Unimplemented", codes.Unimplemented, http.StatusBadRequest},
		{"Internal", codes.Internal, http.StatusTeapot},
		{"Unavailable", codes.Unavailable, http.StatusBadRequest},
		{"DataLoss", codes.DataLoss, http.StatusNotFound},
		{"Default", codes.Code(999), http.StatusTeapot},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPStatusFromCode(tt.code)
			if result != tt.expectedHTTP {
				t.Errorf("HTTPStatusFromCode(%v) = %d, want %d", tt.code, result, tt.expectedHTTP)
			}
		})
	}
}
