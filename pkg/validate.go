package pkg

import (
	"2025_2_404/models"
	"errors"
	"regexp"
)

var allowedPassword = regexp.MustCompile(`^[a-zA-Z0-9._!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$`)
var allowedSymbols = regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
var allowedEmail = regexp.MustCompile(`^[a-zA-Z0-9._]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`);

func validPasswords (password string) bool {
	return allowedPassword.MatchString(password);
}

func validNames (name string) bool {
	return allowedSymbols.MatchString(name)
}

func ValidateRegisterUser (user *models.RegisterUser) error {
	if len(user.UserName)<3 || len(user.UserName)>20{
		return errors.New("Имя пользователя должно быть не меньше 3-х и не больше 20-ти символов")
	}

	if !validNames(user.UserName) {
		return errors.New("Имя пользователя содержит недопустимые значения")
	}

	if !allowedEmail.MatchString(user.Email){
		return errors.New("Недопустимое имя Email")
	}

	if len(user.Password) < 8 {
		return errors.New("Пароль менее 8-ми символов")
	}

	if len(user.Password) > 50 {
		return errors.New("Пароль больше 50-ти символов")
	}

	if !validPasswords(user.Password){
		return errors.New("Недопустимые значения")
	}
	
	return nil
}

func ValidateLoginUser (user *models.BaseUser) error {
	if !allowedEmail.MatchString(user.Email){
		return errors.New("Недопустимое имя Email")
	}

	if len(user.Password) < 8{
		return errors.New("Пароль менее 8-ми символов")
	}

	if len(user.Password) > 50 {
		return errors.New("Пароль больше 50-ти символов")
	}

	if !validPasswords(user.Password){
		return errors.New("Недопустимые значения")
	}
	
	return nil
}