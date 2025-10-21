package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type ID int

type User struct {
	ID          ID	`json:"id"`
	UserName    string	`json:"user_name"`
	Email      string	`json:"email"`
	HashedPassword string	`json:"password"`
}

var allowedSymbols = regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
var allowedEmail = regexp.MustCompile(`^[a-zA-Z0-9._]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`);

func NewUser(userName, email, password string) (*User, error){
	if len(userName)<3 || len(userName)>20{
		return nil, errors.New("username must be at least 3 and no more than 20 characters")
	}

	if !allowedSymbols.MatchString(userName) {
		return nil, errors.New("username contains invalid values")
	}

	if !allowedEmail.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	if len(email) <= 5 || len(email) >= 100 {
		return nil, errors.New("email must be between 5 and 100 characters")
	} 
	
	if len(password) < 8 {
		return nil, errors.New("password less than 8 characters")
	}
	
	if len(password) > 50 {
		return nil, errors.New("password more than 50 characters")
	}
	
	if !allowedSymbols.MatchString(password){
		return nil, errors.New("invalid values")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		UserName:    userName,
		Email:      email,
		HashedPassword: string(hashedPassword),
	}, nil
}

func (u *User) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword),([]byte(password)))
	return err == nil
}

func LoginUser(email, password string) (*User, error){
	if !allowedEmail.MatchString(email){
		return nil, errors.New("invalid email name")
	}
	
	if len(password) < 8{
		return nil, errors.New("password less than 8 characters")
	}
	
	if len(password) > 50 {
		return nil, errors.New("password more than 50 characters")
	}
	
	if !allowedSymbols.MatchString(password){
		return nil, errors.New("invalid values")
	}
	
	return &User{
		Email:      email,
		HashedPassword: password,
	}, nil
}