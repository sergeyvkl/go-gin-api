package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware логирует информацию о каждом запросе
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Старт таймера
		startTime := time.Now()
		
		// Обработка запроса
		c.Next()
		
		// Время окончания обработки
		endTime := time.Now()
		
		// Время выполнения
		latency := endTime.Sub(startTime)
		
		// Метод запроса
		reqMethod := c.Request.Method
		
		// URL запроса
		reqURL := c.Request.URL.Path
		
		// Статус ответа
		statusCode := c.Writer.Status()
		
		// Размер ответа
		resSize := c.Writer.Size()
		
		// Логирование
		log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s | %s | %d bytes",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			c.ClientIP(),
			reqMethod,
			reqURL,
			resSize,
		)
	}
}

// CORSMiddleware добавляет CORS заголовки
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}