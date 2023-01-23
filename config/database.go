package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rizqonmaulana/final-project-edspert/entity"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password= dbname=final_project_edspert sslmode=disable"))

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&entity.Artist{}, &entity.Album{}, &entity.Song{})

	DB = database
}