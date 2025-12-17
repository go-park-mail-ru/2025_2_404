// Package utils provides utility functions for gRPC error handling.
package utils

import (
	"2025_2_404/pkg/globalerrors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, globalerrors.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, globalerrors.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, globalerrors.ErrWrongEmailOrPassword):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, globalerrors.ErrNonValidEmail):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, globalerrors.ErrInvalidQuery):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, globalerrors.ErrNoAuth):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, globalerrors.ErrFileNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, globalerrors.ErrInvalidPath):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, globalerrors.ErrFileWrite),
		errors.Is(err, globalerrors.ErrFileRead),
		errors.Is(err, globalerrors.ErrFileDelete):
		return status.Error(codes.Unknown, err.Error())
	case errors.Is(err, globalerrors.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Unknown, "I'm a teapot")
	}
}
