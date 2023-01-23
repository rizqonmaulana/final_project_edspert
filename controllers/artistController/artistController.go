package artistController

import (
	"net/http"
	"strconv"

	"github.com/rizqonmaulana/final-project-edspert/config"
	"github.com/rizqonmaulana/final-project-edspert/entity"
	helper "github.com/rizqonmaulana/final-project-edspert/helpers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// struct for return get all data
type Artist struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}

func Index(c *gin.Context) {
	var artists []Artist

	// Get query parameters for page number and page size
	pageNumber, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if pageNumber == 0 {
		pageNumber = 1
	}

	if limit == 0 {
		limit = 10
	}

	// Get the total number of artist
	var count int64
	config.DB.Model(&entity.Artist{}).Count(&count)
	
	countInt := int(count)

	// Calculate the offset for the query
	offset := (pageNumber - 1) * limit

	// Execute the paginated query
	config.DB.Select("id","name").Limit(limit).Offset(offset).Find(&artists)

	// Build the URLs for the previous and next pages
	var prevURL, nextURL string
	if pageNumber > 1 {
		prevURL = helper.BuildURL(c, pageNumber-1, limit)
	}
	if (pageNumber * limit) < countInt {
		nextURL = helper.BuildURL(c, pageNumber+1, limit)
	}

	helper.SendResponse(c, http.StatusOK, true, "Success get artists data", nil, gin.H{
		"artists": artists,
		"pagination": gin.H{
			"page": pageNumber,
			"limit": limit,
			"total": count,
			"prev_page_url": prevURL,
			"next_page_url": nextURL,
		},
	})
}

func Show(c *gin.Context) {
	var artist entity.Artist

	id := c.Param("id")

	if err := config.DB.Preload("Albums.Songs").First(&artist, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.SendResponse(c, http.StatusNotFound, false, "Failed to get artist data", []string{"Data not found"}, nil)
			return
		default:
			helper.SendResponse(c, http.StatusInternalServerError, false, "Failed to get artist data", []string{err.Error()}, nil)
			return
		}
	}

	helper.SendResponse(c, http.StatusOK, true, "Success get artist data", nil, gin.H{
		"artist": artist,
	})
}

func Create(c *gin.Context) {
	var artist entity.Artist

	if err := c.ShouldBindJSON(&artist); err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to create artist data", []string{err.Error()}, nil)
		return
	}

	createArtist := config.DB.Create(&artist)

	if createArtist.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to create artist data", []string{createArtist.Error.Error()}, nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, true, "Success created artist data", nil, gin.H{
		"artist": artist,
	})
}

func Update(c *gin.Context) {
	var artist entity.Artist

	id := c.Param("id")

	if err := c.ShouldBindJSON(&artist); err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to update artist data", []string{err.Error()}, nil)
		return
	}

	result := config.DB.Model(&artist).Where("id = ?", id).Updates(&artist)
	
	if result.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to update artist data", []string{result.Error.Error()}, nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, true, "Success updated artist data", nil, nil)
}

func Delete(c *gin.Context) {
	var artist entity.Artist
	id := c.Param("id")

	var albums []entity.Album
	if err := config.DB.Where("artist_id = ?", id).Find(&albums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if config.DB.Where("id = ?", id).Delete(&artist).RowsAffected == 0 {
				helper.SendResponse(c, http.StatusNotFound, false, "Artist not found", []string{"Data not found"}, nil)
				return
			}
			helper.SendResponse(c, http.StatusOK, true, "Success deleted artist data", nil, nil)
			return
		} else {
			helper.SendResponse(c, http.StatusBadRequest, false, "Failed to find albums", []string{err.Error()}, nil)
			return
		}
	}

	for _, album := range albums {
		if err := config.DB.Where("album_id = ?", album.Id).Delete(&entity.Song{}).Error; err != nil {
			helper.SendResponse(c, http.StatusBadRequest, false, "Failed to delete songs", []string{err.Error()}, nil)
			return
		}
	}

	if err := config.DB.Where("artist_id = ?", id).Delete(&entity.Album{}).Error; err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to delete albums", []string{err.Error()}, nil)
		return
	}

	if config.DB.Where("id = ?", id).Delete(&artist).RowsAffected == 0 {
		helper.SendResponse(c, http.StatusNotFound, false, "Artist not found", []string{"Data not found"}, nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, true, "Success deleted artist data", nil, nil)
}







