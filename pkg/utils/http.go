package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.DeadlineExceeded:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusInternalServerError
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.Canceled:
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}
