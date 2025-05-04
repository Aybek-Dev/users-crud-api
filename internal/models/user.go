package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Firstname string    `json:"firstname" binding:"required" gorm:"not null"`
	Lastname  string    `json:"lastname" binding:"required" gorm:"not null"`
	Email     string    `json:"email" binding:"required,email" gorm:"not null;unique"`
	Age       uint      `json:"age" binding:"required,gt=0" gorm:"not null"`
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
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Age       uint   `json:"age" binding:"required,gt=0"`
}

type UserUpdateRequest struct {
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Age       *uint   `json:"age" binding:"omitempty,gt=0"`
}

func (User) TableName() string {
	return "users"
}
