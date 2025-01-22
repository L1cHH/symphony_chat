package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем роутер
	r := gin.Default()

	// Базовый маршрут для проверки
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Запускаем сервер
	log.Println("Starting server at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
