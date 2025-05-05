package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse представляет стандартный формат ответа с ошибкой
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorHandler middleware для централизованной обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Обрабатываем запрос
		c.Next()

		// Если возникли ошибки
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			statusCode := determineStatusCode(c, err)
			errorMsg := formatErrorMessage(err)
			var errorDetails interface{}

			// Проверяем, содержит ли ошибка детали для валидации
			if validationErr, ok := err.(interface {
				ValidationDetails() interface{}
			}); ok {
				errorDetails = validationErr.ValidationDetails()
			} else {
				// Проверка для других типов структурированных ошибок
				if detailedErr, ok := err.(interface {
					Details() interface{}
				}); ok {
					errorDetails = detailedErr.Details()
				}

				// Пытаемся извлечь детали из строки ошибки, если это JSON
				if errorDetails == nil && strings.Contains(errorMsg, "{") && strings.Contains(errorMsg, "}") {
					startIdx := strings.Index(errorMsg, "{")
					endIdx := strings.LastIndex(errorMsg, "}") + 1
					if startIdx >= 0 && endIdx > startIdx {
						jsonStr := errorMsg[startIdx:endIdx]
						var details map[string]interface{}
						if err := json.Unmarshal([]byte(jsonStr), &details); err == nil {
							errorDetails = details
							// Удаляем JSON из сообщения об ошибке
							errorMsg = strings.TrimSpace(errorMsg[:startIdx] + errorMsg[endIdx:])
						}
					}
				}
			}

			// Логируем ошибку
			logErrorWithContext(c, statusCode, errorMsg, errorDetails)

			// Отправляем ответ клиенту
			c.JSON(statusCode, ErrorResponse{
				Status:  statusCode,
				Message: errorMsg,
				Details: errorDetails,
			})
			c.Abort()
		}
	}
}

// determineStatusCode определяет HTTP-статус на основе типа ошибки
func determineStatusCode(c *gin.Context, err error) int {
	// Если статус уже установлен, используем его
	if c.Writer.Status() != http.StatusOK {
		return c.Writer.Status()
	}

	// Проверяем тип ошибки для определения подходящего статуса
	switch err.(type) {
	case *json.SyntaxError, *json.UnmarshalTypeError:
		return http.StatusBadRequest
	default:
		// Проверяем текст ошибки для определения типа
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "не найден"):
			return http.StatusNotFound
		case strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "уже существует"):
			return http.StatusForbidden
		case strings.Contains(errMsg, "validation") || strings.Contains(errMsg, "валидация"):
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}
}

// formatErrorMessage форматирует сообщение об ошибке для ответа пользователю
func formatErrorMessage(err error) string {
	errMsg := err.Error()

	// Для ошибок проверки, возвращаем более дружественное сообщение
	if strings.Contains(errMsg, "validation failed") {
		return "Проверка данных не пройдена"
	}

	// Скрываем внутренние технические детали для пользователя
	if strings.Contains(errMsg, "sql:") {
		return "Ошибка базы данных"
	}

	return errMsg
}

// logErrorWithContext логирует ошибку с контекстом запроса
func logErrorWithContext(c *gin.Context, statusCode int, errorMsg string, details interface{}) {
	// Получаем ID запроса если он есть
	requestID, exists := c.Get("X-Request-ID")
	requestIDStr := "unknown"
	if exists {
		requestIDStr = requestID.(string)
	}

	// Формируем сообщение для лога
	logMsg := "ERROR"
	if statusCode >= 500 {
		logMsg = "CRITICAL ERROR"
	}

	// Детали ошибки
	detailsJson := ""
	if details != nil {
		if detailsBytes, err := json.Marshal(details); err == nil {
			detailsJson = string(detailsBytes)
		}
	}

	// Логируем ошибку
	log.Printf("[%s] RequestID: %s | Path: %s | Status: %d | Message: %s | Details: %s",
		logMsg,
		requestIDStr,
		c.Request.URL.Path,
		statusCode,
		errorMsg,
		detailsJson,
	)
}
