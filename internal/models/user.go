package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Firstname string    `json:"firstname" binding:"required" gorm:"not null" validate:"required,min=2,max=50"`
	Lastname  string    `json:"lastname" binding:"required" gorm:"not null" validate:"required,min=2,max=50"`
	Email     string    `json:"email" binding:"required,email" gorm:"not null;unique" validate:"required,email"`
	Age       uint      `json:"age" binding:"required,gt=0" gorm:"not null" validate:"required,gt=0,lt=120"`
	Created   time.Time `json:"created" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	if u.Created.IsZero() {
		u.Created = time.Now()
	}
	return nil
}

type UserCreateRequest struct {
	Firstname string `json:"firstname" binding:"required" validate:"required,min=2,max=50"`
	Lastname  string `json:"lastname" binding:"required" validate:"required,min=2,max=50"`
	Email     string `json:"email" binding:"required,email" validate:"required,email"`
	Age       uint   `json:"age" binding:"required,gt=0" validate:"required,gt=0,lt=120"`
}

type UserUpdateRequest struct {
	Firstname *string `json:"firstname" validate:"omitempty,min=2,max=50"`
	Lastname  *string `json:"lastname" validate:"omitempty,min=2,max=50"`
	Email     *string `json:"email" binding:"omitempty,email" validate:"omitempty,email"`
	Age       *uint   `json:"age" binding:"omitempty,gt=0" validate:"omitempty,gt=0,lt=120"`
}

type UserResponse struct {
	Success bool   `json:"success"`
	Data    *User  `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UsersResponse struct {
	Success bool   `json:"success"`
	Data    []User `json:"data,omitempty"`
	Count   int64  `json:"count"`
	Error   string `json:"error,omitempty"`
}

func (User) TableName() string {
	return "users"
}
