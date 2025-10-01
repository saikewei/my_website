package photo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(router gin.IRouter) {
	photoGroup := router.Group("/photo")
	{
		photoGroup.POST("/upload", uploadPhoto)
	}
}

func uploadPhoto(c *gin.Context) {
	var newPhoto Photo

	if err := c.ShouldBindJSON(&newPhoto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//TODO: 将照片信息存储到数据库

	c.JSON(http.StatusOK, gin.H{"message": "照片上传成功！", "photo": newPhoto})
}
