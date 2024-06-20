package services

import (
	"tuum.com/internal/models"
	"tuum.com/internal/repositories"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
	return s.userRepository.GetAllUsers()
}
