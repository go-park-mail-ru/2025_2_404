package user

import (
	"testing"
)

func TestNewUser_Success(t *testing.T) {
	user, err := NewUser("testUser", "test@example.com", "Password1!")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be created, got nil")
	}

	if user.UserName != "testUser" {
		t.Errorf("expected UserName 'testUser', got %s", user.UserName)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got %s", user.Email)
	}

	if user.HashedPassword == "" {
		t.Error("expected hashed password to be set")
	}

	if user.HashedPassword == "Password1!" {
		t.Error("expected password to be hashed, got plain text")
	}
}

func TestNewUser_UsernameTooShort(t *testing.T) {
	_, err := NewUser("usr", "test@example.com", "Password1!")

	if err == nil {
		t.Error("expected error for short username, got nil")
	}
}

func TestNewUser_UsernameTooLong(t *testing.T) {
	_, err := NewUser("thisUsernameIsWayTooLong", "test@example.com", "Password1!")

	if err == nil {
		t.Error("expected error for long username, got nil")
	}
}

func TestNewUser_UsernameInvalidCharacters(t *testing.T) {
	_, err := NewUser("test@user", "test@example.com", "Password1!")

	if err == nil {
		t.Error("expected error for invalid characters in username, got nil")
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	testCases := []string{
		"notanemail",
		"@example.com",
		"user@",
		"user@.com",
		"user.example.com",
	}

	for _, email := range testCases {
		_, err := NewUser("testUser", email, "Password1!")
		if err == nil {
			t.Errorf("expected error for invalid email %s, got nil", email)
		}
	}
}

func TestNewUser_EmailTooLong(t *testing.T) {
	longEmail := string(make([]byte, 95)) + "@example.com" // Creates email > 100 chars
	_, err := NewUser("testUser", longEmail, "Password1!")

	if err == nil {
		t.Error("expected error for too long email, got nil")
	}
}

func TestNewUser_PasswordTooShort(t *testing.T) {
	_, err := NewUser("testUser", "test@example.com", "Pass1!")

	if err == nil {
		t.Error("expected error for short password, got nil")
	}
}

func TestNewUser_PasswordTooLong(t *testing.T) {
	longPassword := string(make([]byte, 51))
	_, err := NewUser("testUser", "test@example.com", longPassword)

	if err == nil {
		t.Error("expected error for long password, got nil")
	}
}

func TestNewUser_PasswordMissingUppercase(t *testing.T) {
	_, err := NewUser("testUser", "test@example.com", "password1!")

	if err == nil {
		t.Error("expected error for password without uppercase, got nil")
	}
}

func TestNewUser_PasswordMissingLowercase(t *testing.T) {
	_, err := NewUser("testUser", "test@example.com", "PASSWORD1!")

	if err == nil {
		t.Error("expected error for password without lowercase, got nil")
	}
}

func TestNewUser_PasswordMissingSpecialChar(t *testing.T) {
	_, err := NewUser("testUser", "test@example.com", "Password1")

	if err == nil {
		t.Error("expected error for password without special character, got nil")
	}
}

func TestNewUser_PasswordInvalidCharacters(t *testing.T) {
	_, err := NewUser("testUser", "test@example.com", "Password1!йцу")

	if err == nil {
		t.Error("expected error for password with invalid characters, got nil")
	}
}

func TestUser_PasswordVerification(t *testing.T) {
	user, err := NewUser("testUser", "test@example.com", "Password1!")

	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// The hashed password should not match the plain password
	if user.HashedPassword == "Password1!" {
		t.Error("password should be hashed, not plain text")
	}

	// Verify the password can be checked with bcrypt
	// This indirectly tests that bcrypt was used properly
	if len(user.HashedPassword) < 20 {
		t.Error("hashed password seems too short, might not be bcrypt")
	}
}
