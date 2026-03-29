package database

import (
	"fmt"
	"go-gin-api/config"
	"go-gin-api/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB инициализирует подключение к базе данных
func InitDB(cfg *config.Config) error {
	var err error
	
	dsn := cfg.GetDSN()
	
	// Настройка логгера GORM
	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.Environment == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}
	
	// Подключение к базе данных
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		SkipDefaultTransaction:                   false,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Получение экземпляра sql.DB для настройки пула соединений
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	
	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)           // Максимальное количество простаивающих соединений
	sqlDB.SetMaxOpenConns(100)          // Максимальное количество открытых соединений
	sqlDB.SetConnMaxLifetime(time.Hour) // Максимальное время жизни соединения
	
	log.Println("Database connected successfully")
	return nil
}

// AutoMigrate выполняет автоматическую миграцию моделей
func AutoMigrate() error {
	err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").
		AutoMigrate(
			&models.User{},
			// Добавьте другие модели здесь
	)
	
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	
	log.Println("Database migration completed successfully")
	return nil
}

// CloseDB закрывает подключение к базе данных
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck проверяет состояние подключения к БД
func HealthCheck() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}