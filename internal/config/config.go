package config

import (
	"crud/internal/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Env      string         `yaml:"env"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Config struct {
	DB     *gorm.DB
	Server ServerConfig
}

func LoadConfig() (*AppConfig, error) {
	// Если путь к конфигурации не указан, используем значение по умолчанию
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	// Читаем файл конфигурации
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	// Парсим YAML
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла конфигурации: %w", err)
	}

	// Устанавливаем значения из переменных окружения, если они заданы
	if envHost := os.Getenv("DB_HOST"); envHost != "" {
		config.Database.Host = envHost
	}
	if envPort := os.Getenv("DB_PORT"); envPort != "" {
		config.Database.Port = envPort
	}
	if envUser := os.Getenv("DB_USER"); envUser != "" {
		config.Database.User = envUser
	}
	if envPassword := os.Getenv("DB_PASSWORD"); envPassword != "" {
		config.Database.Password = envPassword
	}
	if envName := os.Getenv("DB_NAME"); envName != "" {
		config.Database.Name = envName
	}
	if envSSLMode := os.Getenv("DB_SSLMODE"); envSSLMode != "" {
		config.Database.SSLMode = envSSLMode
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		config.Server.Port = envPort
	}
	if envEnv := os.Getenv("ENV"); envEnv != "" {
		config.Env = envEnv
	}

	// Преобразуем строки таймаутов в секундах в time.Duration
	config.Server.ReadTimeout = time.Duration(config.Server.ReadTimeout) * time.Second
	config.Server.WriteTimeout = time.Duration(config.Server.WriteTimeout) * time.Second

	return &config, nil
}

func NewConfig() (*Config, error) {
	// Загружаем конфигурацию
	appConfig, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Строка подключения к БД
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.User,
		appConfig.Database.Password,
		appConfig.Database.Name,
		appConfig.Database.SSLMode)

	// Настройка логгера GORM
	gormLogger := logger.Default.LogMode(logger.Info)
	if appConfig.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Открываем соединение с БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Автомиграция схемы
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("ошибка миграции базы данных: %w", err)
	}

	return &Config{
		DB:     db,
		Server: appConfig.Server,
	}, nil
}
