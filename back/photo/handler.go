package photo

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(router gin.IRouter) {
	photoGroup := router.Group("/photo")
	{
		photoGroup.POST("/upload", uploadPhoto)
		photoGroup.POST("/create-album", createAlbum)
	}
}

func uploadPhoto(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未能获取上传的文件"})
		return
	}

	var newPhotoMeta PhotoMeta
	metaStr := c.PostForm("meta")
	if err := json.Unmarshal([]byte(metaStr), &newPhotoMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errChan := make(chan error)

	go func() {
		if err := uploadPhotoFileService(file); err != nil {
			errChan <- err
			return
		}

		//将照片信息存储到数据库
		if err := uploadPhotoMetaStore(newPhotoMeta); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	if err := <-errChan; err != nil {
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

	errChan := make(chan error)

	go func() {
		if err := createAlbumStore(newAlbum); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	if err := <-errChan; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "相册创建成功！", "album": newAlbum})
}
