package usecase

<<<<<<< HEAD:internal/usecase/auth/auth_usecase.go
import ()

=======
import(
	"2025_2_404/internal/domain"
	"errors"
)

type BaseUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	BaseUser
	UserName string `json:"user_name"`
}

type UserUseCase struct {
	userRepo domain.UserRepository
}

// dto - data transport object 

func (uc *UserUseCase) Register(dto RegisterUser) (*domain.User, error) {
	user, err := domain.NewUser(dto.UserName, dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Login(dto BaseUser) (*domain.User, error) {
	user, err := domain.LoginUser(dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}
	if !user.ComparePasswords(dto.Password) {
		return nil, errors.New("Неверный пароль")
	}

	return user, nil
}
>>>>>>> c140dad9a526807051e898f2f7af2f518f8e34ae:internal/usecase/auth_usecase.go
