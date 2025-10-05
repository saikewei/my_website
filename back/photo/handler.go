package photo

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/saikewei/my_website/back/internal/utils"
)

func RegisterRouters(router gin.IRouter) {
	photoGroup := router.Group("/photo")
	{
		photoGroup.POST("/upload", uploadPhoto)
		photoGroup.POST("/create-album", createAlbum)
		photoGroup.POST("/:photo_id/album", addPhotoToAlbum)

		photoGroup.GET("/test", func(c *gin.Context) {
			photo, err := findPhotoByID(6)
			log.Println("photo:", *photo, "Error:", err)

			c.JSON(http.StatusOK, gin.H{"message": "Test endpoint is working!"})
		})
	}
}

func uploadPhoto(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未能获取上传的文件"})
		return
	}

	var newPhotoMeta PhotoUpload
	metaStr := c.PostForm("meta")
	if err := json.Unmarshal([]byte(metaStr), &newPhotoMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utils.RunTaskAsync(func() error {
		if path, size, err := uploadPhotoFileService(file); err != nil {
			return err
		} else if err := uploadPhotoMetaStore(newPhotoMeta, path, size); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应

	c.JSON(http.StatusOK, gin.H{"message": "照片上传成功！", "photo": newPhotoMeta})
}

func createAlbum(c *gin.Context) {
	var newAlbum Album
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := utils.RunTaskAsync(func() error {
		return createAlbumStore(newAlbum)
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "封面照片不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "相册创建成功！", "album": newAlbum})
}

func addPhotoToAlbum(c *gin.Context) {
	var requestBody struct {
		AlbumID int32 `json:"album_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoID := c.Param("photo_id")
	photoIDInt, err := strconv.ParseInt(photoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的照片ID"})
		return
	}

	var pa PhotoAlbum
	pa.PhotoID = int32(photoIDInt)
	pa.AlbumID = requestBody.AlbumID
	err = utils.RunTaskAsync(func() error {
		return addPhotoToAlbumService(pa)
	})

	if err != nil {
		if errors.Is(err, ErrPhotoAlreadyInAlbum) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "照片或相册不存在"}) // 404 Not Found
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"}) // 500
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "照片添加到相册成功！", "photo_album": pa})
	}
}
