package photo

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/saikewei/my_website/back/internal/config"
)

func uploadPhotoFileService(c *gin.Context, file *multipart.FileHeader) error {
	dst := config.C.Storage.PhotoPath + "/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return err
	}
	return nil
}
