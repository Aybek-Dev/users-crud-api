package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"crud/pkg/validator"
	"github.com/gin-gonic/gin"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Message string                    `json:"message"`
	Errors  []validator.ErrorResponse `json:"errors"`
}

// Error возвращает сообщение об ошибке
func (ve *ValidationError) Error() string {
	return ve.Message
}

// ValidationDetails возвращает детали валидации
func (ve *ValidationError) ValidationDetails() interface{} {
	return ve.Errors
}

// ValidateJSON проверяет JSON на соответствие переданной структуре
func ValidateJSON(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что запрос содержит JSON
		if c.Request.Header.Get("Content-Type") != "application/json" {
			err := errors.New("неверный Content-Type, ожидается application/json")
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Считываем тело запроса
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(errors.New("ошибка чтения тела запроса"))
			c.Abort()
			return
		}

		// Восстанавливаем тело запроса
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Привязываем JSON к структуре
		if err := c.ShouldBindJSON(model); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// Валидируем структуру
		validationErrors, err := validator.ValidateStruct(model)
		if err != nil {
			c.Error(&ValidationError{
				Message: "Ошибка валидации данных",
				Errors:  validationErrors,
			})
			c.Abort()
			return
		}

		// Сохраняем валидированную модель в контексте
		c.Set("validated", model)
		c.Next()
	}
}

// GetValidated получает валидированную модель из контекста
func GetValidated(c *gin.Context) interface{} {
	if validated, exists := c.Get("validated"); exists {
		return validated
	}
	return nil
}

// ValidatePathParam проверяет параметр пути
func ValidatePathParam(name string, validationFunc func(string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param(name)
		if param == "" {
			c.Error(errors.New("параметр пути не указан: " + name))
			c.Abort()
			return
		}

		if err := validationFunc(param); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateQueryParam проверяет параметр запроса
func ValidateQueryParam(name string, validationFunc func(string) error, required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Query(name)
		if param == "" && required {
			c.Error(errors.New("обязательный параметр запроса не указан: " + name))
			c.Abort()
			return
		}

		if param != "" {
			if err := validationFunc(param); err != nil {
				c.Error(err)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
