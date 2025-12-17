// Package user defines the User domain model and related types.
package user

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ID = uuid.UUID

type User struct {
	ID             ID     `json:"id"`
	UserName       string `json:"user_name"`
	Email          string `json:"email"`
	HashedPassword string `json:"password"`
	ImagePath      string `json:"img_path"`
	UserFirstName  string `json:"user_first_name"`
	UserLastName   string `json:"user_second_name"`
	Company        string `json:"company"`
	Phone          string `json:"phone_number"`
}

var allowedSymbols = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var allowedPassword = regexp.MustCompile(`^[a-zA-Z0-9._@#$%&+!* =]+$`)
var allowedEmail = regexp.MustCompile(`^[a-zA-Z0-9.+-_]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)
var constUpperCase = regexp.MustCompile(`[A-Z]`)
var constLowerCase = regexp.MustCompile(`[a-z]`)
var constSpecialChar = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)

func NewUser(userName, email, password string) (*User, error) {
	if len(userName) < 4 || len(userName) > 20 {
		return nil, errors.New("username must be at least 4 and no more than 20 characters")
	}

	if !allowedSymbols.MatchString(userName) {
		return nil, errors.New("username contains invalid values")
	}

	if !constLowerCase.MatchString(userName) && !constUpperCase.MatchString(userName) {
		return nil, errors.New("username must contain at least one symbol")
	}

	if !allowedEmail.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	if len(email) >= 100 {
		return nil, errors.New("email must be between 5 and 100 characters")
	}

	if len(password) < 8 || len(password) > 50 {
		return nil, errors.New("password must be between 8 and 50 characters")
	}

	if !allowedPassword.MatchString(password) {
		return nil, errors.New("invalid values")
	}

	if !constLowerCase.MatchString(password) {
		return nil, errors.New("password must contain at least one lower case symbol")
	}

	if !constUpperCase.MatchString(password) {
		return nil, errors.New("password must contain at least one upper case symbol")
	}

	if !constSpecialChar.MatchString(password) {
		return nil, errors.New("password must contain at least one secial symbol")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		UserName:       userName,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}, nil
}

// func LoginUser(email, password string) (*User, error){
// 	if !allowedEmail.MatchString(email){
// 		return nil, errors.New("invalid email name")
// 	}

// 	if len(password) < 8{
// 		return nil, errors.New("password less than 8 characters")
// 	}

// 	if len(password) > 50 {
// 		return nil, errors.New("password more than 50 characters")
// 	}

// 	if !allowedPassword.MatchString(password){
// 		return nil, errors.New("invalid values")
// 	}

// 	return &User{
// 		Email:      email,
// 		HashedPassword: password,
// 	}, nil
// }
