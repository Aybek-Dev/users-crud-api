package tests

import (
	"crud/internal/models"
	"crud/internal/services"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	// Set ID and Created fields to simulate database behavior
	user.ID = uuid.New()
	user.Created = time.Now()
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create service with mock repository
	userService := services.NewUserService(mockRepo)

	// Test data
	createRequest := models.UserCreateRequest{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}

	// Set up expectations
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Call the service method
	user, err := userService.CreateUser(createRequest)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, createRequest.Firstname, user.Firstname)
	assert.Equal(t, createRequest.Lastname, user.Lastname)
	assert.Equal(t, createRequest.Email, user.Email)
	assert.Equal(t, createRequest.Age, user.Age)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.False(t, user.Created.IsZero())

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create service with mock repository
	userService := services.NewUserService(mockRepo)

	// Test data
	userID := uuid.New()
	mockUser := &models.User{
		ID:        userID,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Created:   time.Now(),
	}
	// Set up expectations
	mockRepo.On("GetByID", userID).Return(mockUser, nil)

	// Call the service method
	user, err := userService.GetUser(userID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}
func TestGetUserNotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create service with mock repository
	userService := services.NewUserService(mockRepo)

	// Test data
	userID := uuid.New()

	// Set up expectations
	mockRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	// Call the service method
	user, err := userService.GetUser(userID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create service with mock repository
	userService := services.NewUserService(mockRepo)

	// Test data
	userID := uuid.New()
	mockUser := &models.User{
		ID:        userID,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Created:   time.Now(),
	}

	newFirstname := "Jane"
	updateRequest := models.UserUpdateRequest{
		Firstname: &newFirstname,
	}

	// Set up expectations
	mockRepo.On("GetByID", userID).Return(mockUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Call the service method
	user, err := userService.UpdateUser(userID, updateRequest)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, newFirstname, user.Firstname)

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create service with mock repository
	userService := services.NewUserService(mockRepo)

	// Test data
	userID := uuid.New()

	// Set up expectations
	mockRepo.On("Delete", userID).Return(nil)

	// Call the service method
	err := userService.DeleteUser(userID)

	// Assertions
	assert.NoError(t, err)

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}
