package globalerrors

import (
	"errors"
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
