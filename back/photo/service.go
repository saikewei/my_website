package photo

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/saikewei/my_website/back/internal/config"
	"github.com/saikewei/my_website/back/internal/model"
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

func getOrCreateThumbnailService(photoID int32, originalPath string) (string, error) {
	if originalPath == "" {
		var err error
		originalPath, err = getPhotoPathByIDStore(photoID)
		if err != nil {
			return "", err
		}
	}

	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(filepath.Base(originalPath), ext)
	thumbFileName := base + "_thumb.jpg"
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

	thumb := resize.Resize(400, 0, img, resize.Lanczos3)

	outFile, err := os.Create(thumbPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, thumb, &jpeg.Options{Quality: 75})
	if err != nil {
		os.Remove(thumbPath)
		return "", err
	}

	return thumbPath, nil
}

func getAllPhotosByPageService(page, pageSize int) ([]PhotoAllDataWithThumbnail, int, error) {
	type result struct {
		index int
		data  PhotoAllDataWithThumbnail
		err   error
	}

	photos, total, err := getAllPhotosMetaByPageStore(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(photos) == 0 {
		return []PhotoAllDataWithThumbnail{}, 0, nil
	}

	var wg sync.WaitGroup

	resultsChan := make(chan result, len(photos))
	log.Printf("Starting concurrent processing for %d photos...", len(photos))
	totalLoopStart := time.Now()

	for i, photo := range photos {
		wg.Add(1)
		go func(index int, p *model.VPhotosWithDetail) {
			defer wg.Done()

			thumbPath, err := getOrCreateThumbnailService(photo.ID, p.FilePath)
			if err != nil {
				return
			}

			thumbBase64, err := imageToBase64(thumbPath)
			if err != nil {
				resultsChan <- result{index: index, err: err}
				return
			}
			p.FilePath = "" // 清空原始路径，避免泄露服务器文件结构

			resultsChan <- result{
				index: index,
				data: PhotoAllDataWithThumbnail{
					VPhotosWithDetail: *p,
					ThumbnailBase64:   thumbBase64,
				},
				err: nil,
			}
		}(i, photo)
	}
	wg.Wait()
	close(resultsChan)

	log.Printf("All goroutines finished in %s", time.Since(totalLoopStart))
	finalResults := make([]PhotoAllDataWithThumbnail, len(photos))
	for res := range resultsChan {
		if res.err != nil {
			// 如果任何一个协程失败，我们可以选择立即返回错误
			log.Printf("Error processing photo at index %d: %v", res.index, res.err)
			return nil, 0, res.err
		}
		// 根据原始索引，将结果放入正确的位置
		finalResults[res.index] = res.data
	}

	return finalResults, int(total), nil
}

func imageToBase64(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	encodedString := base64.StdEncoding.EncodeToString(bytes)

	var mimeType string
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	default:
		mimeType = "application/octet-stream" // 默认或未知类型
	}

	return fmt.Sprintf("data:%s;base64,%s", mimeType, encodedString), nil
}
