package app

import (
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindByEmail(email string) (domain.User, error)
	FindById(id uint64) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	Delete(id uint64) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return userService{
		userRepo: userRepository,
	}
}

func (u userService) FindByEmail(email string) (domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u userService) FindById(id uint64) (domain.User, error) {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u userService) Save(user domain.User) (domain.User, error) {
	var err error
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	user, err = u.userRepo.Save(user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u userService) Delete(id uint64) error {
	err := u.userRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
