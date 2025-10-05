package photo

import (
	"errors"
	"io"
	"mime/multipart"
	"os"

	"github.com/saikewei/my_website/back/internal/config"
	"gorm.io/gorm"
)

var ErrPhotoAlreadyInAlbum = errors.New("该照片已经属于一个相册")

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

func addPhotoToAlbumService(pa PhotoAlbum) error {
	if exist, err := checkPhotoExists(pa.PhotoID); err != nil {
		return err
	} else if !exist {
		return gorm.ErrRecordNotFound
	}

	if exist, err := checkAlbumExists(pa.AlbumID); err != nil {
		return err
	} else if !exist {
		return gorm.ErrRecordNotFound
	}

	if modelPhoto, err := findPhotoByID(pa.PhotoID); err != nil {
		return err
	} else if modelPhoto.AlbumID != nil {
		return ErrPhotoAlreadyInAlbum
	} else if err := addPhotoToAlbumStore(pa); err != nil {
		return err
	}

	return nil
}
