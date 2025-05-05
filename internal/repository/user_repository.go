package repository

import (
	"crud/internal/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	GetAll() ([]models.User, int64, error)
	ExistsByEmail(email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *models.User) error {
	start := time.Now()

	result := r.db.Create(user)
	if result.Error != nil {
		log.Printf("DB Error: Не удалось создать пользователя: %v", result.Error)
		return result.Error
	}

	log.Printf("DB: Пользователь создан за %v", time.Since(start))
	return nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	start := time.Now()

	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("DB: Пользователь с ID %s не найден", id)
			return nil, errors.New("user not found")
		}
		log.Printf("DB Error: Не удалось получить пользователя: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("DB: Пользователь получен за %v", time.Since(start))
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	start := time.Now()

	result := r.db.Save(user)
	if result.Error != nil {
		log.Printf("DB Error: Не удалось обновить пользователя: %v", result.Error)
		return result.Error
	}

	log.Printf("DB: Пользователь обновлен за %v", time.Since(start))
	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	start := time.Now()

	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		log.Printf("DB Error: Не удалось удалить пользователя: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("DB: Пользователь с ID %s не найден", id)
		return errors.New("user not found")
	}

	log.Printf("DB: Пользователь удален за %v", time.Since(start))
	return nil
}

func (r *userRepository) GetAll() ([]models.User, int64, error) {
	start := time.Now()

	var users []models.User
	var count int64

	// Получаем общее количество пользователей
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		log.Printf("DB Error: Не удалось получить количество пользователей: %v", err)
		return nil, 0, err
	}

	// Получаем всех пользователей
	if err := r.db.Find(&users).Error; err != nil {
		log.Printf("DB Error: Не удалось получить список пользователей: %v", err)
		return nil, 0, err
	}

	log.Printf("DB: Получено %d пользователей за %v", len(users), time.Since(start))
	return users, count, nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	start := time.Now()

	var count int64
	result := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		log.Printf("DB Error: Не удалось проверить наличие email: %v", result.Error)
		return false, result.Error
	}

	exists := count > 0
	log.Printf("DB: Проверка существования email '%s' (существует: %v) за %v", email, exists, time.Since(start))
	return exists, nil
}
