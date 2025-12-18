package utils

import (
	"2025_2_404/pkg/globalerrors"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPCError_Nil(t *testing.T) {
	result := ToGRPCError(nil)
	if result != nil {
		t.Errorf("ToGRPCError(nil) = %v, want nil", result)
	}
}

func TestToGRPCError_KnownErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode codes.Code
	}{
		{"UserAlreadyExists", globalerrors.ErrUserAlreadyExists, codes.AlreadyExists},
		{"UserNotFound", globalerrors.ErrUserNotFound, codes.NotFound},
		{"WrongEmailOrPassword", globalerrors.ErrWrongEmailOrPassword, codes.Unauthenticated},
		{"NonValidEmail", globalerrors.ErrNonValidEmail, codes.InvalidArgument},
		{"InvalidQuery", globalerrors.ErrInvalidQuery, codes.InvalidArgument},
		{"NoAuth", globalerrors.ErrNoAuth, codes.Unauthenticated},
		{"FileNotFound", globalerrors.ErrFileNotFound, codes.NotFound},
		{"InvalidPath", globalerrors.ErrInvalidPath, codes.InvalidArgument},
		{"FileWrite", globalerrors.ErrFileWrite, codes.Unknown},
		{"FileRead", globalerrors.ErrFileRead, codes.Unknown},
		{"FileDelete", globalerrors.ErrFileDelete, codes.Unknown},
		{"InvalidCredentials", globalerrors.ErrInvalidCredentials, codes.Unauthenticated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGRPCError(tt.err)
			if result == nil {
				t.Fatal("expected error, got nil")
			}

			st, ok := status.FromError(result)
			if !ok {
				t.Fatal("expected gRPC status error")
			}

			if st.Code() != tt.expectedCode {
				t.Errorf("ToGRPCError(%v) code = %v, want %v", tt.err, st.Code(), tt.expectedCode)
			}
		})
	}
}

func TestToGRPCError_UnknownError(t *testing.T) {
	unknownErr := errors.New("some unknown error")
	result := ToGRPCError(unknownErr)

	if result == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(result)
	if !ok {
		t.Fatal("expected gRPC status error")
	}

	if st.Code() != codes.Unknown {
		t.Errorf("ToGRPCError(unknown error) code = %v, want %v", st.Code(), codes.Unknown)
	}

	if st.Message() != "I'm a teapot" {
		t.Errorf("ToGRPCError(unknown error) message = %q, want %q", st.Message(), "I'm a teapot")
	}
}

func TestToGRPCError_WrappedErrors(t *testing.T) {
	wrappedErr := errors.Join(globalerrors.ErrUserNotFound, errors.New("additional context"))
	result := ToGRPCError(wrappedErr)

	if result == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(result)
	if !ok {
		t.Fatal("expected gRPC status error")
	}

	if st.Code() != codes.NotFound {
		t.Errorf("ToGRPCError(wrapped error) code = %v, want %v", st.Code(), codes.NotFound)
	}
}
