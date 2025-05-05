package handlers

import (
	"crud/internal/models"
	"crud/internal/services"
	"crud/pkg/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	//router.POST("/users", h.CreateUser)
	//router.GET("/user/:id", h.GetUser)
	//router.PATCH("/user/:id", h.UpdateUser)
	//router.DELETE("/user/:id", h.DeleteUser)
	router.Use(middleware.ErrorHandler())

	userGroup := router.Group("/api/v1")
	{
		// Создание пользователя
		userGroup.POST("/users", middleware.ValidateJSON(&models.UserCreateRequest{}), h.CreateUser)

		// Получение пользователя по ID
		userGroup.GET("/user/:id", h.validateUUID("id"), h.GetUser)

		// Обновление пользователя
		userGroup.PATCH("/user/:id", h.validateUUID("id"), middleware.ValidateJSON(&models.UserUpdateRequest{}), h.UpdateUser)

		// Удаление пользователя
		userGroup.DELETE("/user/:id", h.validateUUID("id"), h.DeleteUser)

		// Получение всех пользователей (дополнительный эндпоинт)
		userGroup.GET("/users", h.GetAllUsers)
	}
}

func (h *UserHandler) validateUUID(param string) gin.HandlerFunc {
	return middleware.ValidatePathParam(param, func(value string) error {
		_, err := uuid.Parse(value)
		if err != nil {
			return errors.New("неверный формат идентификатора")
		}
		return nil
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	start := time.Now()

	// Получаем валидированные данные из контекста
	requestData, _ := middleware.GetValidated(c).(*models.UserCreateRequest)

	// Создаем пользователя через сервисный слой
	user, err := h.userService.CreateUser(*requestData)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		c.Error(err)
		return
	}

	log.Printf("Пользователь создан успешно, ID: %s, время: %v", user.ID, time.Since(start))

	// Возвращаем созданного пользователя
	c.JSON(http.StatusCreated, models.UserResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// Получаем ID из URL
	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr) // Мы уже проверили ID в middleware

	// Получаем пользователя через сервисный слой
	user, err := h.userService.GetUser(id)
	if err != nil {
		c.Error(err)
		return
	}

	// Возвращаем найденного пользователя
	c.JSON(http.StatusOK, models.UserResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Получаем ID из URL
	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr) // Мы уже проверили ID в middleware

	// Получаем валидированные данные из контекста
	requestData, _ := middleware.GetValidated(c).(*models.UserUpdateRequest)

	// Обновляем пользователя через сервисный слой
	user, err := h.userService.UpdateUser(id, *requestData)
	if err != nil {
		c.Error(err)
		return
	}

	// Возвращаем обновленного пользователя
	c.JSON(http.StatusOK, models.UserResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Получаем ID из URL
	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr) // Мы уже проверили ID в middleware

	// Удаляем пользователя через сервисный слой
	err := h.userService.DeleteUser(id)
	if err != nil {
		c.Error(err)
		return
	}

	// Возвращаем успешный результат
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Пользователь успешно удален",
	})
}
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Получаем пользователей через сервисный слой
	users, count, err := h.userService.GetAllUsers()
	if err != nil {
		c.Error(err)
		return
	}

	// Возвращаем список пользователей
	c.JSON(http.StatusOK, models.UsersResponse{
		Success: true,
		Data:    users,
		Count:   count,
	})
}
