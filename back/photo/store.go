package photo

import (
	"path/filepath"
	"time"

	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/internal/model"
	"github.com/saikewei/my_website/back/internal/model/query"
	"gorm.io/gorm"
)

func uploadPhotoMetaStore(newPhotoMeta PhotoUpload, filePath string, fileSize int64) error {
	var oldTagsID []int32
	var newTags []*model.Tag

	if len(newPhotoMeta.Tags) > 0 {
		for _, tag := range newPhotoMeta.Tags {
			tagID, err := findTagIDByName(tag)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					newTags = append(newTags, &model.Tag{Name: tag})
				} else {
					return err
				}
			} else {
				oldTagsID = append(oldTagsID, tagID)
			}
		}
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		txq := query.Use(tx)

		photo := model.Photo{
			Title:      "无标题",
			FilePath:   filePath,
			FileName:   filepath.Base(filePath),
			FileSize:   int32(fileSize),
			Width:      newPhotoMeta.Width,
			Height:     newPhotoMeta.Height,
			IsFeatured: newPhotoMeta.IsFeatured,
			ShotAt:     newPhotoMeta.ShotAt,
		}
		if err := txq.Photo.Create(&photo); err != nil {
			return err
		}

		metadata := model.PhotoMetadatum{
			PhotoID:      photo.ID,
			Camera:       newPhotoMeta.Camera,
			Lens:         newPhotoMeta.Lens,
			Aperture:     newPhotoMeta.Aperture,
			ShutterSpeed: newPhotoMeta.ShutterSpeed,
			Iso:          newPhotoMeta.Iso,
			ExposureBias: newPhotoMeta.ExposureBias,
			FocalLength:  newPhotoMeta.FocalLength,
			FlashFired:   newPhotoMeta.FlashFired,
			GpsLatitude:  newPhotoMeta.GpsLatitude,
			GpsLongitude: newPhotoMeta.GpsLongitude,
		}
		if err := txq.PhotoMetadatum.Create(&metadata); err != nil {
			return err
		}

		// 处理标签
		if len(newTags) > 0 {
			if err := txq.Tag.CreateInBatches(newTags, 100); err != nil {
				return err
			}

			for _, tag := range newTags {
				oldTagsID = append(oldTagsID, tag.ID)
			}
		}

		if len(oldTagsID) > 0 {
			var photoTags []*model.PhotoTag
			for _, tagID := range oldTagsID {
				photoTags = append(photoTags, &model.PhotoTag{
					PhotoID: photo.ID,
					TagID:   tagID,
				})
			}
			if err := txq.PhotoTag.CreateInBatches(photoTags, 100); err != nil {
				return err
			}
		}

		return nil
	})
}

func findTagIDByName(tagName string) (int32, error) {
	var tag model.Tag
	err := database.DB.Where("name = ?", tagName).First(&tag).Error
	if err != nil {
		return 0, err
	}
	return tag.ID, nil
}

func createAlbumStore(album Album) error {
	modelAlbum := model.Album{
		Title:        album.Title,
		Description:  album.Description,
		CoverPhotoID: album.CoverPhotoID,
	}

	if modelAlbum.CoverPhotoID != nil {
		exist, err := checkPhotoExists(*modelAlbum.CoverPhotoID)
		if err != nil {
			return err
		} else if !exist {
			return gorm.ErrRecordNotFound
		}
	}

	return database.DB.Create(&modelAlbum).Error
}

func checkPhotoExists(photoID int32) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Photo{}).Where("id = ?", photoID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func checkAlbumExists(albumID int32) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Album{}).Where("id = ?", albumID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func findPhotoByID(photoID int32) (*model.Photo, error) {
	var photo model.Photo
	err := database.DB.First(&photo, photoID).Error
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func addPhotoToAlbumStore(pa PhotoAlbum) error {
	result := database.DB.Model(&model.Photo{}).
		Where("id = ?", pa.PhotoID).
		Update("album_id", pa.AlbumID)

	// 2. 检查执行过程中是否发生错误
	if result.Error != nil {
		return result.Error
	}

	// 3. 检查是否有行被实际更新
	if result.RowsAffected == 0 {
		// 如果没有行被更新，说明传入的 PhotoID 不存在
		return gorm.ErrRecordNotFound
	}

	return nil
}

func getAllAlbumsIDStore() ([]int32, error) {
	var albumIDs []int32
	err := database.DB.Model(&model.Album{}).Order("created_at asc").Pluck("id", &albumIDs).Error
	if err != nil {
		return nil, err
	}
	return albumIDs, nil
}

func getAlbumByIDStore(albumID int32) (*Album, error) {
	var album model.Album
	err := database.DB.First(&album, albumID).Error
	if err != nil {
		return nil, err
	}
	return &Album{
		ID:           album.ID,
		Title:        album.Title,
		Description:  album.Description,
		CoverPhotoID: album.CoverPhotoID,
		CreatedAt:    album.CreatedAt,
		UpdatedAt:    album.UpdatedAt,
	}, nil
}

func getPhotoPathByIDStore(photoID int32) (string, error) {
	var photo model.Photo
	err := database.DB.Select("file_path").First(&photo, photoID).Error
	if err != nil {
		return "", err
	}
	return photo.FilePath, nil
}

func editAlbumStore(album *Album) error {
	modelAlbum := model.Album{
		Title:       album.Title,
		Description: album.Description,
		UpdatedAt:   time.Now(),
	}

	return database.DB.Model(&model.Album{}).Where("id = ?", album.ID).Updates(&modelAlbum).Error
}
