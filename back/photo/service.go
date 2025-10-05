package photo

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/saikewei/my_website/back/internal/config"
	"gorm.io/gorm"
)

var ErrPhotoAlreadyInAlbum = errors.New("该照片已经属于一个相册")

func uploadPhotoFileService(file *multipart.FileHeader) (string, int64, error) {
	src, err := file.Open()
	if err != nil {
		return "", 0, err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)

	newFileName := uuid.NewString() + ext

	// 创建目标文件
	dstPath := filepath.Join(config.C.Storage.PhotoPath, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", 0, err
	}
	defer dst.Close()

	// 将上传文件的内容拷贝到目标文件
	if _, err = io.Copy(dst, src); err != nil {
		return "", 0, err
	}

	return dstPath, file.Size, nil
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
