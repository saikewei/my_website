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

func RegisterRouters(publicGroup, privateGroup gin.IRouter) {
	publicPhotoGroup := publicGroup.Group("/photo")
	{
		publicPhotoGroup.GET("/albums-id", getAllAlbumsID)
		publicPhotoGroup.GET("/album/details", getAllAlbumsDetails)
		publicPhotoGroup.GET("/:photo_id", getPhotoByID)
		publicPhotoGroup.GET("/:photo_id/thumbnail", getPhotoThumbnailByID)
		publicPhotoGroup.GET("/album/:album_id", getAlbumByID)
		publicPhotoGroup.GET("/page", getPhotosByPage)
		publicPhotoGroup.GET("/photos", getPhotosByCursor)

		publicPhotoGroup.GET("/test", func(c *gin.Context) {
			_, total, _ := getAllPhotosMetaByPageStore(1, 5)
			log.Printf("Total photos: %d", total)

			c.JSON(http.StatusOK, gin.H{"message": "Test endpoint is working!"})
		})
	}

	privatePhotoGroup := privateGroup.Group("/photo")
	{
		privatePhotoGroup.POST("/upload", uploadPhoto)
		privatePhotoGroup.POST("/create-album", createAlbum)
		privatePhotoGroup.POST("/:photo_id/album", addPhotoToAlbum)

		privatePhotoGroup.PUT("/edit/album", editAlbum)
		privatePhotoGroup.PUT("/edit/album-cover", editAlbumCover)
		privatePhotoGroup.PUT("/edit/photo", editPhoto)

		privatePhotoGroup.DELETE("/album/:album_id", deleteAlbum)
		privatePhotoGroup.DELETE("/:photo_id", deletePhoto)
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
		if _, err := getOrCreateThumbnailService(newPhotoMeta.ID, "", 0); err != nil {
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

func getAllAlbumsDetails(c *gin.Context) {
	albums, err := getAllAlbumsDetailsStore()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"albums": albums})
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
	thumbnailSize := c.DefaultQuery("size", "400")
	sizeInt, err := strconv.ParseUint(thumbnailSize, 10, 32)
	if err != nil || sizeInt < 100 || sizeInt > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的缩略图尺寸，必须在100到2000之间"})
		return
	}

	photoID := c.Param("photo_id")
	photoIDInt, err := strconv.ParseInt(photoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的照片ID"})
		return
	}

	// 调用服务层函数来获取或创建缩略图
	thumbPath, err := getOrCreateThumbnailService(int32(photoIDInt), "", uint(sizeInt))

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

func deleteAlbum(c *gin.Context) {
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

func getPhotosByPage(c *gin.Context) {
	pageNumStr := c.DefaultQuery("page-num", "1")
	pageSizeStr := c.DefaultQuery("page-size", "10")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页码"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 5 || pageSize > 40 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的每页大小，必须在5到40之间"})
		return
	}

	photosWithThumbnails, total, err := getAllPhotosByPageService(int(pageNum), int(pageSize))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "没有找到照片"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photos": photosWithThumbnails, "total": total})
}

func editPhoto(c *gin.Context) {
	newPhoto := PhotoEdit{}
	if err := c.ShouldBindJSON(&newPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newPhoto.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "照片ID不能为空"})
		return
	}

	err := editPhotoService(newPhoto)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "照片不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "照片编辑成功！", "photo": newPhoto})
}

func deletePhoto(c *gin.Context) {
	photoID := c.Param("photo_id")
	photoIDInt, err := strconv.ParseInt(photoID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的照片ID"})
		return
	}

	err = deletePhotoByIDStore(int32(photoIDInt))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "照片不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "照片删除成功！"})
}

func getPhotosByCursor(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 5 || limit > 40 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 limit 参数，必须在5到40之间"})
		return
	}

	cursorStr := c.DefaultQuery("cursor", "0")
	cursor, err := strconv.ParseInt(cursorStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 cursor 参数"})
		return
	}

	albumStr := c.DefaultQuery("album", "0")
	album, err := strconv.ParseInt(albumStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 album 参数"})
		return
	}

	photos, err := getPhotosByCursorService(int32(cursor), limit, int32(album))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor int32
	hasMore := false
	if len(photos) > 0 {
		// 如果获取到的数据量等于请求量，我们认为可能还有更多数据
		if len(photos) == limit {
			hasMore = true
			// 将最后一个元素的ID作为下一次请求的游标
			nextCursor = photos[len(photos)-1].ID
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"photos":      photos,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

func editAlbumCover(c *gin.Context) {
	newCover := AlbumCover{}
	if err := c.ShouldBindJSON(&newCover); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newCover.AlbumID == 0 || newCover.CoverPhotoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为0"})
	}

	err := editAlbumCoverStore(newCover)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "相册不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "相册封面设置成功！"})
}
