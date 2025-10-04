package photo

import (
	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/internal/model"
	"github.com/saikewei/my_website/back/internal/model/query"
	"gorm.io/gorm"
)

func uploadPhotoMetaStore(newPhotoMeta PhotoMeta) error {
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
			AlbumID:     newPhotoMeta.AlbumID,
			Title:       newPhotoMeta.Title,
			Description: newPhotoMeta.Description,
			FilePath:    newPhotoMeta.FilePath,
			FileName:    newPhotoMeta.FileName,
			FileSize:    newPhotoMeta.FileSize,
			Width:       newPhotoMeta.Width,
			Height:      newPhotoMeta.Height,
			IsFeatured:  newPhotoMeta.IsFeatured,
			ShotAt:      newPhotoMeta.ShotAt,
		}
		if err := txq.Photo.Create(&photo); err != nil {
			return err
		}

		metadata := model.PhotoMetadatum{
			PhotoID:      photo.ID,
			CameraID:     newPhotoMeta.CameraID,
			LensID:       newPhotoMeta.LensID,
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
