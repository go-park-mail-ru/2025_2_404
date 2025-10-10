package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type ID int64

type User struct {
	ID          ID
	UserName    string
	Email      string
	HashedPassword string
}

var allowedSymbols = regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
var allowedEmail = regexp.MustCompile(`^[a-zA-Z0-9._]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`);

func NewUser(userName, email, password string) (*User, error){
	if len(userName)<3 || len(userName)>20{
		return nil, errors.New("имя пользователя должно быть не меньше 3-х и не больше 20-ти символов")
	}

	if !allowedSymbols.MatchString(userName) {
		return nil, errors.New("имя пользователя содержит недопустимые значения")
	}

	if !allowedEmail.MatchString(email) {
		return nil, errors.New("недопустимое имя Email")
	}
	
	if len(password) < 8 {
		return nil, errors.New("пароль менее 8-ми символов")
	}
	
	if len(password) > 50 {
		return nil, errors.New("пароль больше 50-ти символов")
	}
	
	if !allowedSymbols.MatchString(password){
		return nil, errors.New("недопустимые значения")
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
		return nil, errors.New("недопустимое имя Email")
	}
	
	if len(password) < 8{
		return nil, errors.New("пароль менее 8-ми символов")
	}
	
	if len(password) > 50 {
		return nil, errors.New("пароль больше 50-ти символов")
	}
	
	if !allowedSymbols.MatchString(password){
		return nil, errors.New("недопустимые значения")
	}
	
	return &User{
		Email:      email,
		HashedPassword: password,
	}, nil
}