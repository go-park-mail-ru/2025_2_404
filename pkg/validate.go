package pkg

import (
	"2025_2_404/models"
	"errors"
	"regexp"
	"strings"
)

func validPasswords (password string) bool {
	forbidenSymb := " :'\"";
	
	if strings.ContainsAny(password, forbidenSymb){
		return false 
	}
	return true;
}

func validNames (name string) bool {
	forbidenSymb := "!@#$%^&*()[]{};:'\"<>,.?/\\|`~ ";

	if strings.ContainsAny(name, forbidenSymb) {
		return false
	}
	return true
}

var allowedEmail = regexp.MustCompile("^[a-Za-Z0-9._]+@[a-Za-Z0-9-]+.[a-zA-Z]{2,}$");

func ValidateRegisterUser (user *models.RegisterUser) error {
	if len(user.UserName)<3 {
		return errors.New("Имя пользователя должно быть не меньше 3-х символов")
	}

	if !validNames(user.UserName) {
		return errors.New("Имя пользователя содержит недопустимые значения")
	}

	if !allowedEmail.MatchString(user.Email){
		return errors.New("Недопустимое имя Email")
	}

	if len(user.Password) < 8{
		return errors.New("Пароль менее 8-ми символов")
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

	if !validPasswords(user.Password){
		return errors.New("Недопустимые значения")
	}
	
	return nil
}