package service

import (
	modeluser "2025_2_404/internal/service/auth/domain"
	"2025_2_404/pkg/globalerrors"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type mockRepository struct {
	createFunc      func(ctx context.Context, user *modeluser.User) (modeluser.ID, error)
	findByEmailFunc func(ctx context.Context, email string) (modeluser.User, error)
}

func (m *mockRepository) Create(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return uuid.Nil, nil
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (modeluser.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, email)
	}
	return modeluser.User{}, globalerrors.ErrUserNotFound
}

type mockTokenUsecase struct {
	generateTokenFunc func(userID modeluser.ID) (string, error)
}

func (m *mockTokenUsecase) GenerateToken(userID modeluser.ID) (string, error) {
	if m.generateTokenFunc != nil {
		return m.generateTokenFunc(userID)
	}
	return "", nil
}

func TestRegister_Success(t *testing.T) {
	expectedID := uuid.New()
	expectedToken := "test-token-123"

	mockRepo := &mockRepository{
		createFunc: func(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
			return expectedID, nil
		},
	}

	mockToken := &mockTokenUsecase{
		generateTokenFunc: func(userID modeluser.ID) (string, error) {
			if userID == expectedID {
				return expectedToken, nil
			}
			return "", errors.New("wrong user ID")
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	token, userID, err := useCase.Register(context.Background(), "test@example.com", "Password1!", "validUser")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if token != expectedToken {
		t.Errorf("expected token %s, got %s", expectedToken, token)
	}

	if userID != expectedID {
		t.Errorf("expected userID %v, got %v", expectedID, userID)
	}
}

func TestRegister_ValidationError(t *testing.T) {
	mockRepo := &mockRepository{}
	mockToken := &mockTokenUsecase{}
	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	_, _, err := useCase.Register(context.Background(), "bad-email", "Password1!", "validUser")
	if err == nil {
		t.Error("expected validation error, got nil")
	}
}

func TestRegister_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	mockRepo := &mockRepository{
		createFunc: func(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
			return uuid.Nil, expectedErr
		},
	}

	mockToken := &mockTokenUsecase{}
	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	_, _, err := useCase.Register(context.Background(), "test@example.com", "Password1!", "validUser")
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestRegister_TokenGenerationError(t *testing.T) {
	expectedErr := errors.New("token generation failed")

	mockRepo := &mockRepository{
		createFunc: func(ctx context.Context, user *modeluser.User) (modeluser.ID, error) {
			return uuid.New(), nil
		},
	}

	mockToken := &mockTokenUsecase{
		generateTokenFunc: func(userID modeluser.ID) (string, error) {
			return "", expectedErr
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	_, _, err := useCase.Register(context.Background(), "test@example.com", "Password1!", "validUser")
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestCheck_Success(t *testing.T) {
	expectedID := uuid.New()
	password := "Password1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo := &mockRepository{
		findByEmailFunc: func(ctx context.Context, email string) (modeluser.User, error) {
			return modeluser.User{
				ID:             expectedID,
				Email:          email,
				HashedPassword: string(hashedPassword),
			}, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, nil, logger)

	userID, err := useCase.Check(context.Background(), "test@example.com", password)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if userID != expectedID {
		t.Errorf("expected userID %v, got %v", expectedID, userID)
	}
}

func TestCheck_UserNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		findByEmailFunc: func(ctx context.Context, email string) (modeluser.User, error) {
			return modeluser.User{}, globalerrors.ErrUserNotFound
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, nil, logger)

	_, err := useCase.Check(context.Background(), "test@example.com", "Password1!")
	if err != globalerrors.ErrWrongEmailOrPassword {
		t.Errorf("expected ErrWrongEmailOrPassword, got %v", err)
	}
}

func TestCheck_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("RightPassword1!"), bcrypt.DefaultCost)

	mockRepo := &mockRepository{
		findByEmailFunc: func(ctx context.Context, email string) (modeluser.User, error) {
			return modeluser.User{
				ID:             uuid.New(),
				Email:          email,
				HashedPassword: string(hashedPassword),
			}, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, nil, logger)

	_, err := useCase.Check(context.Background(), "test@example.com", "WrongPassword1!")
	if err != globalerrors.ErrWrongEmailOrPassword {
		t.Errorf("expected ErrWrongEmailOrPassword, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	expectedID := uuid.New()
	expectedToken := "test-token-123"
	password := "Password1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo := &mockRepository{
		findByEmailFunc: func(ctx context.Context, email string) (modeluser.User, error) {
			return modeluser.User{
				ID:             expectedID,
				Email:          email,
				HashedPassword: string(hashedPassword),
			}, nil
		},
	}

	mockToken := &mockTokenUsecase{
		generateTokenFunc: func(userID modeluser.ID) (string, error) {
			return expectedToken, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	token, userID, err := useCase.Login(context.Background(), "test@example.com", password)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if token != expectedToken {
		t.Errorf("expected token %s, got %s", expectedToken, token)
	}

	if userID != expectedID {
		t.Errorf("expected userID %v, got %v", expectedID, userID)
	}
}

func TestLogin_ValidationError(t *testing.T) {
	mockRepo := &mockRepository{}
	mockToken := &mockTokenUsecase{}
	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	_, _, err := useCase.Login(context.Background(), "bad-email", "Password1!")
	if err != globalerrors.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_CheckFailed(t *testing.T) {
	mockRepo := &mockRepository{
		findByEmailFunc: func(ctx context.Context, email string) (modeluser.User, error) {
			return modeluser.User{}, globalerrors.ErrUserNotFound
		},
	}

	mockToken := &mockTokenUsecase{}
	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	_, _, err := useCase.Login(context.Background(), "test@example.com", "Password1!")
	if err != globalerrors.ErrWrongEmailOrPassword {
		t.Errorf("expected ErrWrongEmailOrPassword, got %v", err)
	}
}

func TestLogin_TokenGenerationError(t *testing.T) {
	password := "Password1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedErr := errors.New("token generation failed")

	mockRepo := &mockRepository{
		findByEmailFunc: func(ctx context.Context, email string) (modeluser.User, error) {
			return modeluser.User{
				ID:             uuid.New(),
				Email:          email,
				HashedPassword: string(hashedPassword),
			}, nil
		},
	}

	mockToken := &mockTokenUsecase{
		generateTokenFunc: func(userID modeluser.ID) (string, error) {
			return "", expectedErr
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, mockToken, logger)

	_, _, err := useCase.Login(context.Background(), "test@example.com", password)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestNew_UseCase(t *testing.T) {
	mockRepo := &mockRepository{}
	mockToken := &mockTokenUsecase{}
	logger, _ := zap.NewDevelopment()

	useCase := New(mockRepo, mockToken, logger)

	if useCase == nil {
		t.Fatal("expected non-nil useCase")
	}

	if useCase.repo == nil {
		t.Error("expected non-nil repo")
	}

	if useCase.tokenUsecase == nil {
		t.Error("expected non-nil tokenUsecase")
	}

	if useCase.logger == nil {
		t.Error("expected non-nil logger")
	}
}
