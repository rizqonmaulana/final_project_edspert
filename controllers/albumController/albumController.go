package albumController

import (
	"net/http"
	"strconv"

	"github.com/rizqonmaulana/final-project-edspert/config"
	"github.com/rizqonmaulana/final-project-edspert/entity"
	helper "github.com/rizqonmaulana/final-project-edspert/helpers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Album struct {
	Id int64 `json:"id"`
	Artist_Id int64 `json:"artist_id"`
	Title string `json:"title"`
	Price int64 `json:"price"`
}

func Index(c *gin.Context) {
	var albums []Album

	db := config.DB

	// Get query parameters for page number and page size
	pageNumber, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if artistID := c.Query("artist_id"); artistID != "" {
		db = db.Where("artist_id = ?", artistID)
	}

	if pageNumber == 0 {
		pageNumber = 1
	}

	if limit == 0 {
		limit = 10
	}

	// Get the total number of albums
	var count int64
	db.Model(&entity.Album{}).Count(&count)
	
	countInt := int(count)

	// Calculate the offset for the query
	offset := (pageNumber - 1) * limit

	// Execute the paginated query
	db.Limit(limit).Offset(offset).Find(&albums)

	// Build the URLs for the previous and next pages
	var prevURL, nextURL string
	if pageNumber > 1 {
		prevURL = helper.BuildURL(c, pageNumber-1, limit)
	}
	if (pageNumber * limit) < countInt {
		nextURL = helper.BuildURL(c, pageNumber+1, limit)
	}

	helper.SendResponse(c, http.StatusOK, true, "Success get albums data", nil, gin.H{
		"albums": albums,
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
	var album entity.Album

	id := c.Param("id")

	if err := config.DB.Preload("Songs").First(&album, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.SendResponse(c, http.StatusNotFound, false, "Failed to get album data", []string{"Data not found"}, nil)
			return
		default:
			helper.SendResponse(c, http.StatusInternalServerError, false, "Failed to get album data", []string{err.Error()}, nil)
			return
		}
	}

	helper.SendResponse(c, http.StatusOK, true, "Success get album data", nil, gin.H{
		"album": album,
	})
}

func Create(c *gin.Context) {
	var album entity.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to create album data", []string{err.Error()}, nil)
		return
	}

	createAlbum := config.DB.Create(&album)

	if createAlbum.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to create album data", []string{createAlbum.Error.Error()}, nil)
		return
	}
	
	helper.SendResponse(c, http.StatusOK, true, "Success created album data", nil, gin.H{
		"album": album,
	})
}

func Update(c *gin.Context) {
	var album entity.Album

	id := c.Param("id")

	if err := c.ShouldBindJSON(&album); err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to update album data", []string{err.Error()}, nil)
		return
	}

	result := config.DB.Model(&album).Where("id = ?", id).Updates(&album)
	if result.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to update album data", []string{result.Error.Error()}, nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, true, "Success updated album data", nil, nil)
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	config.DB.Where("album_id = ?", id).Delete(&entity.Song{})
	result := config.DB.Where("id = ?", id).Delete(&entity.Album{})

	if result.RowsAffected == 0 {
		helper.SendResponse(c, http.StatusNotFound, false, "Album not found", []string{"Data not found"}, nil)
		return
	}
	
	if result.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to delete album data", []string{result.Error.Error()}, nil)
		return
	}
	
	helper.SendResponse(c, http.StatusOK, true, "Success deleted album data", nil, nil)
}