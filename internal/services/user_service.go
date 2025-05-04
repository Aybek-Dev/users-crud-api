package services

import (
	"crud/internal/models"
	"crud/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(request models.UserCreateRequest) (*models.User, error)
	GetUser(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, request models.UserUpdateRequest) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(request models.UserCreateRequest) (*models.User, error) {
	user := &models.User{
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Email:     request.Email,
		Age:       request.Age,
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateUser(id uuid.UUID, request models.UserUpdateRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates only for fields that are provided
	if request.Firstname != nil {
		user.Firstname = *request.Firstname
	}
	if request.Lastname != nil {
		user.Lastname = *request.Lastname
	}
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Age != nil {
		user.Age = *request.Age
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}
