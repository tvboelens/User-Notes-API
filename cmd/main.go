package main

import (
	"user-notes-api/config"
	"user-notes-api/models"
	"user-notes-api/routes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" +
		cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	db.AutoMigrate(&models.User{}, &models.Note{})

	r := gin.Default()
	routes.SetupRoutes(r, db, cfg.JWTSecret)
	r.Run(":" + cfg.AppPort)

}
