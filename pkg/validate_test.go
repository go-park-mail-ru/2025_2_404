package pkg

import (
	"2025_2_404/internal/models"
	"testing"
	"github.com/stretchr/testify/assert"
)

// Тестируем конструктор NewUser, который содержит логику валидации для регистрации
func TestNewUser(t *testing.T) {
	testCases := []struct {
		name          string
		userName      string
		email         string
		password      string
		expectError   bool
		expectedError string
	}{
		{
			name:        "Успешное создание пользователя",
			userName:    "ValidUser",
			email:       "test.user@example.com",
			password:    "Password123",
			expectError: false,
		},
		{
			name:          "Слишком короткое имя пользователя",
			userName:      "me",
			email:         "test.user@example.com",
			password:      "Password123",
			expectError:   true,
			expectedError: "Имя пользователя должно быть не меньше 3-х и не больше 20-ти символов",
		},
		{
			name:          "Недопустимые символы в имени пользователя",
			userName:      "!hehe!",
			email:         "test.user@example.com",
			password:      "Password123",
			expectError:   true,
			expectedError: "Имя пользователя содержит недопустимые значения",
		},
		{
			name:          "Недопустимый формат Email",
			userName:      "ValidUser",
			email:         "testuser@example.com@",
			password:      "Password123",
			expectError:   true,
			expectedError: "Недопустимое имя Email",
		},
		{
			name:          "Слишком короткий пароль",
			userName:      "ValidUser",
			email:         "testuser@example.com",
			password:      "Pass12",
			expectError:   true,
			expectedError: "Пароль менее 8-ми символов",
		},
		{
			name:          "Слишком длинный пароль",
			userName:      "ValidUser",
			email:         "testuser@example.com",
			password:      "Wv!8kz#rEpA6q@T3sB*uY4n^C5hF@dG2jL!mK9vP*sX7zQc#gH1", // > 50 символов
			expectError:   true,
			expectedError: "Пароль больше 50-ти символов",
		},
		{
			name:          "Недопустимые символы в пароле (слеш)",
			userName:      "ValidUser",
			email:         "testuser@example.com",
			password:      "Pass:/123",
			expectError:   true,
			expectedError: "Недопустимые значения",
		},
		{
			name:          "Недопустимые символы в пароле (кириллица)",
			userName:      "ValidUser",
			email:         "testuser@example.com",
			password:      "пароль123",
			expectError:   true,
			expectedError: "Недопустимые значения",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := models.NewUser(tc.userName, tc.email, tc.password)
			if tc.expectError {
				assert.Error(t, err) // Проверяем, что ошибка есть
				assert.Nil(t, user)  // Убеждаемся, что пользователь не создан
				assert.Equal(t, tc.expectedError, err.Error()) // Сравниваем текст ошибки
			} else {
				assert.NoError(t, err) // Проверяем, что ошибки нет
				assert.NotNil(t, user) // Убеждаемся, что пользователь создан
				assert.Equal(t, tc.userName, user.UserName)
				assert.Equal(t, tc.email, user.Email)
			}
		})
	}
}

// Тестируем конструктор LoginUser, который содержит логику валидации для входа
func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		password      string
		expectError   bool
		expectedError string
	}{
		{
			name:        "Успешный вход",
			email:       "test@example.com",
			password:    "ValidPassword123",
			expectError: false,
		},
		{
			name:          "Недопустимый email (кириллица)",
			email:         "ярусский@gmail.com",
			password:      "123Passw3",
			expectError:   true,
			expectedError: "Недопустимое имя Email",
		},
		{
			name:          "Пароль менее 8 символов",
			email:         "test@gmail.com",
			password:      "123P3",
			expectError:   true,
			expectedError: "Пароль менее 8-ми символов",
		},
		{
			name:          "Пароль более 50 символов",
			email:         "test@gmail.com",
			password:      "Wv!8kz#rEpA6q@T3sB*uY4n^C5hF@dG2jL!mK9vP*sX7zQc#gH1",
			expectError:   true,
			expectedError: "Пароль больше 50-ти символов",
		},
		{
			name:          "Недопустимые символы в пароле",
			email:         "test@gmail.com",
			password:      "1672 :ж",
			expectError:   true,
			expectedError: "Недопустимые значения",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := models.LoginUser(tc.email, tc.password)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}