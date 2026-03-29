package routes

import (
	"go-gin-api/controllers"
	"go-gin-api/middleware"
	"go-gin-api/repositories"
	"go-gin-api/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Инициализация репозитория, сервиса и контроллера
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	
	// Добавляем middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	
	// Health check endpoint
	router.GET("/health", userController.HealthCheck)
	
	// Группа API маршрутов
	api := router.Group("/api/v1")
	{
		// Маршруты для пользователей
		users := api.Group("/users")
		{
			users.GET("", userController.GetAllUsers)
			users.GET("/:id", userController.GetUserByID)
			users.POST("", userController.CreateUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
		
		// Можно добавить дополнительные группы
		// api.Group("/products")
		// api.Group("/orders")
	}
	
	// Маршрут для несуществующих путей
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not found",
			"message": "The requested resource was not found",
		})
	})
}