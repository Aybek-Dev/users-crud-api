package services

import (
	"crud/internal/models"
	"crud/internal/repository"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type UserService interface {
	CreateUser(request models.UserCreateRequest) (*models.User, error)
	GetUser(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, request models.UserUpdateRequest) (*models.User, error)
	DeleteUser(id uuid.UUID) error
	GetAllUsers() ([]models.User, int64, error)
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
	start := time.Now()

	// Проверяем, существует ли пользователь с таким email
	exists, err := s.userRepo.ExistsByEmail(request.Email)
	if err != nil {
		log.Printf("Ошибка при проверке email: %v", err)
		return nil, errors.New("ошибка при создании пользователя")
	}

	if exists {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Создаем модель пользователя
	user := &models.User{
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Email:     request.Email,
		Age:       request.Age,
	}

	// Сохраняем пользователя в БД
	err = s.userRepo.Create(user)
	if err != nil {
		log.Printf("Ошибка при создании пользователя в БД: %v", err)
		return nil, errors.New("ошибка при создании пользователя")
	}

	log.Printf("Пользователь создан успешно за %v", time.Since(start))

	return user, nil
}

func (s *userService) GetUser(id uuid.UUID) (*models.User, error) {
	start := time.Now()

	// Получаем пользователя из БД
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		if err.Error() == "user not found" {
			return nil, errors.New("пользователь не найден")
		}
		return nil, errors.New("ошибка при получении пользователя")
	}

	log.Printf("Пользователь получен успешно за %v", time.Since(start))

	return user, nil
}

func (s *userService) UpdateUser(id uuid.UUID, request models.UserUpdateRequest) (*models.User, error) {
	start := time.Now()

	// Получаем пользователя из БД
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Ошибка при получении пользователя для обновления: %v", err)
		if err.Error() == "user not found" {
			return nil, errors.New("пользователь не найден")
		}
		return nil, errors.New("ошибка при обновлении пользователя")
	}

	// Проверяем, не пытаемся ли изменить email на уже существующий
	if request.Email != nil && *request.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(*request.Email)
		if err != nil {
			log.Printf("Ошибка при проверке email: %v", err)
			return nil, errors.New("ошибка при обновлении пользователя")
		}

		if exists {
			return nil, errors.New("пользователь с таким email уже существует")
		}
	}

	// Применяем обновления только для указанных полей
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

	// Сохраняем обновленного пользователя
	err = s.userRepo.Update(user)
	if err != nil {
		log.Printf("Ошибка при обновлении пользователя в БД: %v", err)
		return nil, errors.New("ошибка при обновлении пользователя")
	}

	log.Printf("Пользователь обновлен успешно за %v", time.Since(start))

	return user, nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	start := time.Now()

	// Удаляем пользователя из БД
	err := s.userRepo.Delete(id)
	if err != nil {
		log.Printf("Ошибка при удалении пользователя: %v", err)
		if err.Error() == "user not found" {
			return errors.New("пользователь не найден")
		}
		return errors.New("ошибка при удалении пользователя")
	}

	log.Printf("Пользователь удален успешно за %v", time.Since(start))

	return nil
}

func (s *userService) GetAllUsers() ([]models.User, int64, error) {
	start := time.Now()

	// Получаем пользователей из БД
	users, count, err := s.userRepo.GetAll()
	if err != nil {
		log.Printf("Ошибка при получении списка пользователей: %v", err)
		return nil, 0, errors.New("ошибка при получении списка пользователей")
	}

	log.Printf("Получен список пользователей (всего: %d) за %v", count, time.Since(start))

	return users, count, nil
}
