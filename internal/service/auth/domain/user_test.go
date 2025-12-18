package domain

import (
	"2025_2_404/pkg/globalerrors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidateRegisterUser_Success(t *testing.T) {
	user, err := ValidateRegisterUser("validUser", "test@example.com", "Password1!")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.UserName != "validUser" {
		t.Errorf("expected username 'validUser', got %s", user.UserName)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte("Password1!"))
	if err != nil {
		t.Errorf("password hash is invalid: %v", err)
	}
}

func TestValidateRegisterUser_UsernameTooShort(t *testing.T) {
	_, err := ValidateRegisterUser("usr", "test@example.com", "Password1!")
	if err != globalerrors.ErrUsernameTooShort {
		t.Errorf("expected ErrUsernameTooShort, got %v", err)
	}
}

func TestValidateRegisterUser_UsernameTooLong(t *testing.T) {
	_, err := ValidateRegisterUser("verylongusernameover20", "test@example.com", "Password1!")
	if err != globalerrors.ErrUsernameTooLong {
		t.Errorf("expected ErrUsernameTooLong, got %v", err)
	}
}

func TestValidateRegisterUser_UsernameInvalidChars(t *testing.T) {
	_, err := ValidateRegisterUser("user@name", "test@example.com", "Password1!")
	if err != globalerrors.ErrUsernameInvalidChars {
		t.Errorf("expected ErrUsernameInvalidChars, got %v", err)
	}
}

func TestValidateRegisterUser_UsernameNoLetters(t *testing.T) {
	_, err := ValidateRegisterUser("12345", "test@example.com", "Password1!")
	if err != globalerrors.ErrUsernameNoLetters {
		t.Errorf("expected ErrUsernameNoLetters, got %v", err)
	}
}

func TestValidateRegisterUser_InvalidEmail(t *testing.T) {
	testCases := []struct {
		name  string
		email string
	}{
		{"missing @", "testexample.com"},
		{"missing domain", "test@"},
		{"missing TLD", "test@example"},
		{"invalid chars", "test@@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateRegisterUser("validUser", tc.email, "Password1!")
			if err != globalerrors.ErrNonValidEmail {
				t.Errorf("expected ErrNonValidEmail for %s, got %v", tc.email, err)
			}
		})
	}
}

func TestValidateRegisterUser_EmailTooLong(t *testing.T) {
	longEmail := "verylongemailaddressthatexceedsthemaximumlengthallowedforthistestingpurposeandmoreandmore@exampledomainname.com"
	if len(longEmail) < 100 {
		t.Skipf("email is not long enough: %d chars", len(longEmail))
	}
	_, err := ValidateRegisterUser("validUser", longEmail, "Password1!")
	if err != globalerrors.ErrNonValidEmail {
		t.Errorf("expected ErrNonValidEmail for long email, got %v", err)
	}
}

func TestValidateRegisterUser_PasswordTooShort(t *testing.T) {
	_, err := ValidateRegisterUser("validUser", "test@example.com", "Pass1!")
	if err != globalerrors.ErrPasswordTooShort {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestValidateRegisterUser_PasswordTooLong(t *testing.T) {
	longPassword := "VeryLongPassword1!VeryLongPassword1!VeryLongPassword1!"
	_, err := ValidateRegisterUser("validUser", "test@example.com", longPassword)
	if err != globalerrors.ErrPasswordTooLong {
		t.Errorf("expected ErrPasswordTooLong, got %v", err)
	}
}

func TestValidateRegisterUser_PasswordInvalidChars(t *testing.T) {
	_, err := ValidateRegisterUser("validUser", "test@example.com", "Password1!😀")
	if err != globalerrors.ErrPasswordInvalidChars {
		t.Errorf("expected ErrPasswordInvalidChars, got %v", err)
	}
}

func TestValidateRegisterUser_PasswordNoLower(t *testing.T) {
	_, err := ValidateRegisterUser("validUser", "test@example.com", "PASSWORD1!")
	if err != globalerrors.ErrPasswordNoLower {
		t.Errorf("expected ErrPasswordNoLower, got %v", err)
	}
}

func TestValidateRegisterUser_PasswordNoUpper(t *testing.T) {
	_, err := ValidateRegisterUser("validUser", "test@example.com", "password1!")
	if err != globalerrors.ErrPasswordNoUpper {
		t.Errorf("expected ErrPasswordNoUpper, got %v", err)
	}
}

func TestValidateRegisterUser_PasswordNoSpecial(t *testing.T) {
	_, err := ValidateRegisterUser("validUser", "test@example.com", "Password1")
	if err != globalerrors.ErrPasswordNoSpecial {
		t.Errorf("expected ErrPasswordNoSpecial, got %v", err)
	}
}

func TestValidateLoginUser_Success(t *testing.T) {
	err := ValidateLoginUser("test@example.com", "Password1!")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateLoginUser_EmptyEmail(t *testing.T) {
	err := ValidateLoginUser("", "Password1!")
	if err != globalerrors.ErrEmailRequired {
		t.Errorf("expected ErrEmailRequired, got %v", err)
	}
}

func TestValidateLoginUser_InvalidEmail(t *testing.T) {
	err := ValidateLoginUser("notanemail", "Password1!")
	if err != globalerrors.ErrNonValidEmail {
		t.Errorf("expected ErrNonValidEmail, got %v", err)
	}
}

func TestValidateLoginUser_EmptyPassword(t *testing.T) {
	err := ValidateLoginUser("test@example.com", "")
	if err != globalerrors.ErrPasswordRequired {
		t.Errorf("expected ErrPasswordRequired, got %v", err)
	}
}

func TestValidateLoginUser_PasswordTooShort(t *testing.T) {
	err := ValidateLoginUser("test@example.com", "Pass1!")
	if err != globalerrors.ErrPasswordTooShort {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestValidateLoginUser_PasswordNoLower(t *testing.T) {
	err := ValidateLoginUser("test@example.com", "PASSWORD1!")
	if err != globalerrors.ErrPasswordNoLower {
		t.Errorf("expected ErrPasswordNoLower, got %v", err)
	}
}

func TestValidateLoginUser_PasswordNoUpper(t *testing.T) {
	err := ValidateLoginUser("test@example.com", "password1!")
	if err != globalerrors.ErrPasswordNoUpper {
		t.Errorf("expected ErrPasswordNoUpper, got %v", err)
	}
}

func TestValidateLoginUser_PasswordNoSpecial(t *testing.T) {
	err := ValidateLoginUser("test@example.com", "Password1")
	if err != globalerrors.ErrPasswordNoSpecial {
		t.Errorf("expected ErrPasswordNoSpecial, got %v", err)
	}
}
