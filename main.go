package main

import (
	"go-gin-api/config"
	"go-gin-api/database"
	"go-gin-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	
	// Инициализация базы данных
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	
	// Выполнение миграций
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	// Создание роутера
	router := gin.Default()
	
	// Настройка маршрутов
	routes.SetupRoutes(router)
	
	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}