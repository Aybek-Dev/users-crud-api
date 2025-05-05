package main

import (
	"context"
	"crud/internal/config"
	"crud/internal/handlers"
	"crud/internal/repository"
	"crud/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	log.Println("Запуск приложения...")
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Ошибка при получении текущей директории: %v", err)
	} else {
		log.Printf("Текущая рабочая директория: %s", currentDir)
	}

	// Загружаем конфигурацию
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Установка режима Gin
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	// Initialize router
	router := gin.Default()

	// Add middleware for request logging, CORS, etc.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(cfg.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers and register routes
	userHandler := handlers.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server is running on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
