package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rizqonmaulana/final-project-edspert/entity"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("POSTGRES_URL")
	database, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&entity.Artist{}, &entity.Album{}, &entity.Song{})

	DB = database
}