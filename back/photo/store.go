package photo

import (
	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/internal/model"
	"github.com/saikewei/my_website/back/internal/model/query"
	"gorm.io/gorm"
)

func uploadPhotoMetaStore(newPhotoMeta PhotoMeta) error {

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

		if len(newPhotoMeta.TagsID) > 0 {
			var photoTags []*model.PhotoTag
			for _, tagID := range newPhotoMeta.TagsID {
				photoTags = append(photoTags, &model.PhotoTag{
					PhotoID: photo.ID,
					TagID:   tagID,
				})
			}
			if err := txq.PhotoTag.CreateInBatches(photoTags, len(photoTags)); err != nil {
				return err
			}
		}

		return nil
	})
}
