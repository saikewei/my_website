package photo

import (
	"errors"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
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

func getOrCreateThumbnailService(photoID int32) (string, error) {
	originalPath, err := getPhotoPathByIDStore(photoID)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(filepath.Base(originalPath), ext)
	thumbFileName := base + "_thumb" + ext
	thumbPath := filepath.Join(filepath.Dir(originalPath), thumbFileName)

	if _, err := os.Stat(thumbPath); err == nil {
		return thumbPath, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	originalFile, err := os.Open(originalPath)
	if err != nil {
		return "", err
	}
	defer originalFile.Close()

	img, _, err := image.Decode(originalFile)
	if err != nil {
		return "", err
	}

	thumb := resize.Resize(1000, 0, img, resize.Lanczos3)

	outFile, err := os.Create(thumbPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	err = png.Encode(outFile, thumb)
	if err != nil {
		return "", err
	}

	return thumbPath, nil
}
