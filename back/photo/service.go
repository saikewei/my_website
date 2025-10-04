package photo

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/saikewei/my_website/back/internal/config"
)

func uploadPhotoFileService(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dstPath := config.C.Storage.PhotoPath + "/" + file.Filename
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 将上传文件的内容拷贝到目标文件
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
