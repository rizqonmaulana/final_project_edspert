package songController

import (
	"net/http"
	"strconv"

	"github.com/rizqonmaulana/final-project-edspert/config"
	"github.com/rizqonmaulana/final-project-edspert/entity"
	helper "github.com/rizqonmaulana/final-project-edspert/helpers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
    var songs []entity.Song
    var albums []entity.Album

    db := config.DB

    // Get query parameters for page number, page size, and artist ID
    pageNumber, _ := strconv.Atoi(c.Query("page"))
    limit, _ := strconv.Atoi(c.Query("limit"))
	albumID := c.Query("album_id")
    artistID := c.Query("artist_id")

	if artistID != "" {
		db.Where("artist_id = ?", artistID).Find(&albums)
	}
	if albumID != "" {
		db = db.Where("album_id = ?", albumID)
	}
	

    if pageNumber == 0 {
        pageNumber = 1
    }

    if limit == 0 {
        limit = 10
    }

    // Get the total number of songs
    var count int64
    db.Model(&entity.Song{}).Count(&count)

    countInt := int(count)

    // Calculate the offset for the query
    offset := (pageNumber - 1) * limit

    if len(albums) > 0 {
        albumIDs := make([]uint, len(albums))
        for i, a := range albums {
            albumIDs[i] = uint(a.Id)
        }

        // Execute the paginated query with the album IDs as a filter
        db.Where("album_id IN (?)", albumIDs).Limit(limit).Offset(offset).Find(&songs)
    } else {
        // Execute the paginated query without any album ID filter
        db.Limit(limit).Offset(offset).Find(&songs)
    }

    // Build the URLs for the previous and next pages
    var prevURL, nextURL string
    if pageNumber > 1 {
        prevURL = helper.BuildURL(c, pageNumber-1, limit)
    }
    if (pageNumber * limit) < countInt {
        nextURL = helper.BuildURL(c, pageNumber+1, limit)
    }

    helper.SendResponse(c, http.StatusOK, true, "Success get songs data", nil, gin.H{
        "songs": songs,
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
	var song entity.Song

	id := c.Param("id")

	if err := config.DB.First(&song, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.SendResponse(c, http.StatusNotFound, false, "Failed to get song data", []string{"Data not found"}, nil)
			return
		default:
			helper.SendResponse(c, http.StatusInternalServerError, false, "Failed to get song data", []string{err.Error()}, nil)
			return
		}
	}

	helper.SendResponse(c, http.StatusOK, true, "Success get song data", nil, gin.H{
		"song": song,
	})
}

func Create(c *gin.Context) {
	var song entity.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to create song data", []string{err.Error()}, nil)
		return
	}

	createSong := config.DB.Create(&song)

	if createSong.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to create song data", []string{createSong.Error.Error()}, nil)
		return
	}
	
	helper.SendResponse(c, http.StatusOK, true, "Success created song data", nil, gin.H{
		"song": song,
	})
}

func Update(c *gin.Context) {
	var song entity.Song

	id := c.Param("id")

	if err := c.ShouldBindJSON(&song); err != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to update song data", []string{err.Error()}, nil)
		return
	}

	result := config.DB.Model(&song).Where("id = ?", id).Updates(&song)
	
	if result.Error != nil {
		helper.SendResponse(c, http.StatusBadRequest, false, "Failed to update song data", []string{result.Error.Error()}, nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, true, "Success updated song data", nil, nil)
}

func Delete(c *gin.Context) {
	var song entity.Song
	
	id := c.Param("id")

	if config.DB.Delete(&song, id).RowsAffected == 0 {
		helper.SendResponse(c, http.StatusNotFound, false, "Song not found", []string{"Data not found"}, nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, true, "Success deleted song data", nil, nil)
}