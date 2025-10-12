package main

import (
	"user-notes-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := config.LoadConfig()
}
