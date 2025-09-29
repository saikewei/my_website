package photo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(router *gin.Engine) {
	// photoGroup := router.Group("/photos")
	// {
	// 	photoGroup.GET("/", getPhotos)
	// 	photoGroup.POST("/", uploadPhoto)
	// }
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is the photo handler!",
		})
	})
}
