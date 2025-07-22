package main

import (
	"UserSegmentationServise/internal/handlers"
	"UserSegmentationServise/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// DSN: строка подключения к PostgreSQL
	dsn := "host=localhost user=postgres password=postgres dbname=segments port=5433 sslmode=disable"

	// Открытие соединения с БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Миграция схемы БД при запуске
	db.AutoMigrate(&models.Segment{}, &models.UserSegment{}, &models.SegmentRule{})

	// Инициализация роутера Gin
	r := gin.Default()

	// Группа роутов с префиксом /api
	api := r.Group("/api")
	handlers.RegisterRoutes(api, db)

	// Запуск сервера на порту 8080
	r.Run(":8080")
}
