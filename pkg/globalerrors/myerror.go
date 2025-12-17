package globalerrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidQuery          = errors.New("invalid query parameter")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrWrongEmailOrPassword  = errors.New("wrong email or password")
	ErrNonValidEmail         = errors.New("invalid email")
	ErrInternal              = errors.New("internal error")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong       = errors.New("password must be no more than 50 characters")
	ErrPasswordNoLower       = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoUpper       = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoSpecial     = errors.New("password must contain at least one special character")
	ErrPasswordInvalidChars  = errors.New("password contains invalid characters")
	ErrUsernameTooShort      = errors.New("username must be at least 4 characters")
	ErrUsernameTooLong       = errors.New("username must be no more than 20 characters")
	ErrUsernameInvalidChars  = errors.New("username contains invalid characters")
	ErrUsernameNoLetters     = errors.New("username must contain at least one letter")
	ErrEmailRequired         = errors.New("email is required")
	ErrPasswordRequired      = errors.New("password is required")
	ErrorContextTimeout		 = errors.New("dedline timeout")
	ErrNoAuth    			 = errors.New("no auth")
)

func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case err == ErrUserAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case err == ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case err == ErrWrongEmailOrPassword:
		return status.Error(codes.Unauthenticated, err.Error())
	case err == ErrNonValidEmail:
		return status.Error(codes.InvalidArgument, err.Error())
	case err == ErrInvalidQuery:
		return status.Error(codes.InvalidArgument, err.Error())
	case err == ErrNoAuth:
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
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
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}