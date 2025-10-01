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
			user: &models.RegisterUser{ 
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
			expectedError: "Имя пользователя должно быть не меньше 3-х и не больше 20-ти символов",
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
				BaseUser: models.BaseUser{Email: "testuser@example.com@", Password: "Password123"}, 
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
			nameTest: "Пароль больше 50-ти символов",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "testuser@example.com", Password: "Wv!8k$z#rE&pA6q@T3sB*uY4n^C5hF@dG2jL!mK9vP*sX7zQc#gH1"},
				UserName: "meow",
			},
			expectError:   true,
			expectedError: "Пароль больше 50-ти символов",
		},
		{
			nameTest: "Недопустимые значения пароля",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "testuser@example.com", Password: "pasdf1560~"},
				UserName: "meow",
			},
			expectError:   true,
			expectedError: "Недопустимые значения",
		},
		{
			nameTest: "Недопустимые значения пароля",
			user: &models.RegisterUser{
				BaseUser: models.BaseUser{Email: "testuser@example.com", Password: "пароль12"},
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
		user *models.BaseUser	
		expectError bool
		expectedError string
	}{
		{
			nameTest: "unКорректный email",
			user: &models.BaseUser{Email: "ярусский@gmail.com", Password: "123Passw3"},
			expectError: false,
			expectedError: "Недопустимое имя Email",
		},
		{
			nameTest: "Пароль менее 8 символов",
			user: &models.BaseUser{Email: "test@gmail.com", Password: "123P3"},
			expectError: false,
			expectedError: "Пароль менее 8-ми символов",
		},
		{
			nameTest: "Пароль более 8 символов",
			user: &models.BaseUser{Email: "test@gmail.com", Password: "Wv!8k$z#rE&pA6q@T3sB*uY4n^C5hF@dG2jL!mK9vP*sX7zQc#gH1"},
			expectError: false,
			expectedError: "Пароль больше 50-ти символов",
		},
		{
			nameTest: "Недопустимые значения пароля",
			user: &models.BaseUser{Email: "test@gmail.com", Password: "1672 :ж"},
			expectError: false,
			expectedError: "Недопустимые значения",
		},
	}

	for _, tc := range testCases{
		t.Run(tc.nameTest, func(t *testing.T) {
			res := ValidateLoginUser(tc.user)
			assert.Equal(t, tc.expectedError, res.Error())
		})
	}
}