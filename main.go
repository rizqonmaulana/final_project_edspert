package main

import (
	"github.com/gin-gonic/gin"

	"github.com/rizqonmaulana/final-project-edspert/config"
	"github.com/rizqonmaulana/final-project-edspert/controllers/albumController"
	"github.com/rizqonmaulana/final-project-edspert/controllers/artistController"
	"github.com/rizqonmaulana/final-project-edspert/controllers/songController"
)

func main() {
	r := gin.Default();
	config.ConnectDatabase()

	// artist route
	r.GET("/api/artists", artistController.Index)
	r.GET("/api/artist/:id", artistController.Show)
	r.POST("/api/artist", artistController.Create)
	r.PUT("/api/artist/:id", artistController.Update)
	r.DELETE("/api/artist/:id", artistController.Delete)

	// album route
	r.GET("/api/albums", albumController.Index)
	r.GET("/api/album/:id", albumController.Show)
	r.POST("/api/album", albumController.Create)
	r.PUT("/api/album/:id", albumController.Update)
	r.DELETE("/api/album/:id", albumController.Delete)
	
	// song route
	r.GET("/api/songs", songController.Index)
	r.GET("/api/song/:id", songController.Show)
	r.POST("/api/song", songController.Create)
	r.PUT("/api/song/:id", songController.Update)
	r.DELETE("/api/song/:id", songController.Delete)

	r.Run()
}