package pkg

import (
	"2025_2_404/models"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestValidateRegisterUser(t *testing.T){
	testCases := []struct{
		nameTest string
		user *models.RegisterUser
		expectError bool
		expectedError string
	}{
		{
			nameTest: "Успешная валидация", 
			user: &models.RegisterUser{ // Используем & и именованные поля
				BaseUser: models.BaseUser{Email: "test.user@example.com", Password: "Password123"},
				UserName: "ValidUser",
			},
			expectError: false,
		},
		{
			nameTest: "Короткое имя пользователя!",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "test.user@example.com", Password: "Password123"},
				UserName: "me",
			},
			expectError:   true,
			expectedError: "Имя пользователя должно быть не меньше 3-х символов",
		},
		{
			nameTest: "Недопустимые значения в имени",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "test.user@example.com", Password: "Password123"},
				UserName: "!hehe!",
			},
			expectError:   true,
			expectedError: "Имя пользователя содержит недопустимые значения",
		},
		{
			nameTest: "Недопустимое имя Email",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "testuser@example.com@", Password: "Password123"}, // Исправил email для теста
				UserName: "meow",
			},
			expectError:   true,
			expectedError: "Недопустимое имя Email",
		},
		{
			nameTest: "Пароль менее 8-ми символов",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "testuser@example.com", Password: "Pass12"},
				UserName: "meow",
			},
			expectError:   true,
			expectedError: "Пароль менее 8-ми символов",
		},
		{
			nameTest: "Недопустимые значения пароля",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "testuser@example.com", Password: "Pass:/12"},
				UserName: "meow",
			},
			expectError:   true,
			expectedError: "Недопустимые значения",
		},
	}

	for _, tc := range testCases{
		t.Run(tc.nameTest, func(t *testing.T) {
			err := ValidateRegisterUser(tc.user);
			if tc.expectError {
				assert.Error(t, err) 
				assert.Equal(t, tc.expectedError, err.Error()) 
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLoginUser(t *testing.T){
	testCases := []struct{
		nameTest string
		password string
		expectedError bool
	}{
		{
			nameTest: "Корректный пароль",
			password: "Password123-ok",
			expectedError: true,
		},
		{
			nameTest: "Недопустимый пароль с пробелом",
			password: "Password 123",
			expectedError: false,
		},
		{
			nameTest: "Недопустимый пароль с двоеточием",
			password: "Password:123",
			expectedError: false,
		},
		{
			nameTest: "Недопустимый пароль с одинарной кавычкой",
			password: "Password'123", 
			expectedError: false,
		},
		{
			nameTest: "Недопустимый пароль с двойными кавычками",
			password: "Password\"123",
			expectedError: false,
		},
		{
			nameTest: "Пустой пароль (валиден по текущей логике)",
			password: "",
			expectedError: true, 
		},
	}

	for _, tc := range testCases{
		t.Run(tc.nameTest, func(t *testing.T) {
			res := validPasswords(tc.password)
			assert.Equal(t, tc.expectedError, res)
		})
	}
}