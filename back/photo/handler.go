package photo

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRouters(router gin.IRouter) {
	photoGroup := router.Group("/photo")
	{
		photoGroup.POST("/upload", uploadPhoto)
		photoGroup.POST("/create-album", createAlbum)
		photoGroup.POST("/:photo_id/album", addPhotoToAlbum)

		photoGroup.GET("/albums-id", getAllAlbumsID)
		photoGroup.GET("/:photo_id", getPhotoByID)
		photoGroup.GET("/:photo_id/thumbnail", getPhotoThumbnailByID)
		photoGroup.GET("/album/:album_id", getAlbumByID)

		photoGroup.GET("/test", func(c *gin.Context) {
			photo, err := findPhotoByID(6)
			log.Println("photo:", *photo, "Error:", err)

			c.JSON(http.StatusOK, gin.H{"message": "Test endpoint is working!"})
		})

		photoGroup.PUT("/edit-album", editAlbum)

		photoGroup.DELETE("/album/:album_id", deleteAlbumByID)
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

	if path, size, err := uploadPhotoFileService(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if newPhotoID, err := uploadPhotoMetaStore(newPhotoMeta, path, size); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		newPhotoMeta.ID = newPhotoID
	}

	go func() {
		if _, err := getOrCreateThumbnailService(newPhotoMeta.ID); err != nil {
			log.Println("生成缩略图失败:", err)
		} else {
			log.Println("缩略图生成成功")
		}
	}()
	// 返回成功响应

	c.JSON(http.StatusOK, gin.H{"message": "照片上传成功！", "photo": newPhotoMeta})
}

func createAlbum(c *gin.Context) {
	var newAlbum Album
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newAlbum.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "相册标题不能为空"})
		return
	}

	err := createAlbumStore(newAlbum)

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
	err = addPhotoToAlbumService(pa)

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

func getAllAlbumsID(c *gin.Context) {
	albumsID, err := getAllAlbumsIDStore()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"albums_id": albumsID})
}

func getAlbumByID(c *gin.Context) {
	albumID := c.Param("album_id")
	albumIDInt, err := strconv.ParseInt(albumID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	album, err := getAlbumByIDStore(int32(albumIDInt))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"album": album})
}

func getPhotoByID(c *gin.Context) {
	photoID := c.Param("photo_id")
	photoIDInt, err := strconv.ParseInt(photoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的照片ID"})
		return
	}

	photo, err := getPhotoPathByIDStore(int32(photoIDInt))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "照片不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.File(photo)
}

func getPhotoThumbnailByID(c *gin.Context) {
	photoID := c.Param("photo_id")
	photoIDInt, err := strconv.ParseInt(photoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的照片ID"})
		return
	}

	// 调用服务层函数来获取或创建缩略图
	thumbPath, err := getOrCreateThumbnailService(int32(photoIDInt))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "照片不存在"})
		} else {
			// 其他错误，例如文件读写、图片解码等
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 将文件发送给客户端
	c.File(thumbPath)
}

func editAlbum(c *gin.Context) {
	newAlbum := Album{}
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newAlbum.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "相册ID不能为空"})
		return
	}
	if newAlbum.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "相册标题不能为空"})
		return
	}

	err := editAlbumStore(&newAlbum)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "相册不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "相册编辑成功！", "album": newAlbum})
}

func deleteAlbumByID(c *gin.Context) {
	albumID := c.Param("album_id")
	albumIDInt, err := strconv.ParseInt(albumID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的相册ID"})
		return
	}

	err = deleteAlbumByIDStore(int32(albumIDInt))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "相册删除成功！"})
}
