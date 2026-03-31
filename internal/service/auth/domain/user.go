// Package domain defines the User domain model and related types.
package domain

import (
	"2025_2_404/pkg/globalerrors"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ID = uuid.UUID

type User struct {
	ID             ID
	UserName       string
	Email          string
	HashedPassword string
}

var allowedSymbols = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var allowedPassword = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=\[\]{};':"\\|,.<>?/-]+$`)
var allowedEmail = regexp.MustCompile(`^[a-zA-Z0-9._+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var constUpperCase = regexp.MustCompile(`[A-Z]`)
var constLowerCase = regexp.MustCompile(`[a-z]`)
var constSpecialChar = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)

func ValidateRegisterUser(userName, email, password string) (*User, error) {
	if len(userName) < 4 {
		return nil, globalerrors.ErrUsernameTooShort
	}
	if len(userName) > 20 {
		return nil, globalerrors.ErrUsernameTooLong
	}

	if !allowedSymbols.MatchString(userName) {
		return nil, globalerrors.ErrUsernameInvalidChars
	}

	if !constLowerCase.MatchString(userName) && !constUpperCase.MatchString(userName) {
		return nil, globalerrors.ErrUsernameNoLetters
	}

	if !allowedEmail.MatchString(email) {
		return nil, globalerrors.ErrNonValidEmail
	}

	if len(email) >= 100 {
		return nil, globalerrors.ErrNonValidEmail // или отдельная ошибка, но обычно длина покрывается email-валидацией
	}

	if len(password) < 8 {
		return nil, globalerrors.ErrPasswordTooShort
	}
	if len(password) > 50 {
		return nil, globalerrors.ErrPasswordTooLong
	}

	if !allowedPassword.MatchString(password) {
		return nil, globalerrors.ErrPasswordInvalidChars
	}

	if !constLowerCase.MatchString(password) {
		return nil, globalerrors.ErrPasswordNoLower
	}

	if !constUpperCase.MatchString(password) {
		return nil, globalerrors.ErrPasswordNoUpper
	}

	if !constSpecialChar.MatchString(password) {
		return nil, globalerrors.ErrPasswordNoSpecial
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, globalerrors.ErrInternal
	}

	return &User{
		ID:             uuid.New(),
		UserName:       userName,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}, nil
}

func ValidateLoginUser(email, password string) error {
	if email == "" {
		return globalerrors.ErrEmailRequired
	}

	if !allowedEmail.MatchString(email) {
		return globalerrors.ErrNonValidEmail
	}

	if password == "" {
		return globalerrors.ErrPasswordRequired
	}

	if len(password) < 8 {
		return globalerrors.ErrPasswordTooShort
	}
	if len(password) > 50 {
		return globalerrors.ErrPasswordTooLong
	}

	if !allowedPassword.MatchString(password) {
		return globalerrors.ErrPasswordInvalidChars
	}

	if !constLowerCase.MatchString(password) {
		return globalerrors.ErrPasswordNoLower
	}

	if !constUpperCase.MatchString(password) {
		return globalerrors.ErrPasswordNoUpper
	}

	if !constSpecialChar.MatchString(password) {
		return globalerrors.ErrPasswordNoSpecial
	}

	return nil
}
