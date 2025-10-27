package photo

import (
	"path/filepath"
	"time"

	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/internal/model"
	"github.com/saikewei/my_website/back/internal/model/query"
	"gorm.io/gorm"
)

func uploadPhotoMetaStore(newPhotoMeta PhotoUpload, filePath string, fileSize int64) (int32, error) {
	var oldTagsID []int32
	var newTags []*model.Tag

	if len(newPhotoMeta.Tags) > 0 {
		for _, tag := range newPhotoMeta.Tags {
			tagID, err := findTagIDByName(tag)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					newTags = append(newTags, &model.Tag{Name: tag})
				} else {
					return 0, err
				}
			} else {
				oldTagsID = append(oldTagsID, tagID)
			}
		}
	}
	var newPhotoID int32

	err := database.DB.Transaction(func(tx *gorm.DB) error {
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
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := txq.Photo.Create(&photo); err != nil {
			return err
		}

		newPhotoID = photo.ID

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

	if err != nil {
		return 0, err
	}

	return newPhotoID, nil
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

func getAllAlbumsDetailsStore() ([]*Album, error) {
	var albums []model.Album
	err := database.DB.Model(&model.Album{}).Order("created_at asc").Find(&albums).Error
	if err != nil {
		return nil, err
	}

	var result []*Album
	for _, album := range albums {
		result = append(result, &Album{
			ID:           album.ID,
			Title:        album.Title,
			Description:  album.Description,
			CoverPhotoID: album.CoverPhotoID,
			CreatedAt:    album.CreatedAt,
			UpdatedAt:    album.UpdatedAt,
		})
	}

	return result, nil
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

func deleteAlbumByIDStore(albumID int32) error {
	result := database.DB.Model(&model.Album{}).Where("id = ?", albumID).Delete(&model.Album{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func getAllPhotosMetaByPageStore(page, pageSize int) ([]*model.VPhotosWithDetail, int64, error) {
	var photos []*model.VPhotosWithDetail
	var total int64

	err := database.DB.Model(&model.Photo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	} else if (page-1)*pageSize >= int(total) {
		return nil, total, gorm.ErrRecordNotFound
	}

	err = database.DB.Model(&model.VPhotosWithDetail{}).Order("shot_at desc").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&photos).Error
	if err != nil {
		return nil, 0, err
	}

	return photos, total, nil
}

func getAllPhotoPageNumStore(pageSize int) (int, error) {
	var total int64
	err := database.DB.Model(&model.Photo{}).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return int((total + int64(pageSize) - 1) / int64(pageSize)), nil
}

func editPhotoByIDStore(newData PhotoEdit) error {
	oldPhoto, err := findPhotoByID(newData.ID)
	if err != nil {
		return err
	}

	var changeAlbum bool
	if oldPhoto.AlbumID == nil {
		if *newData.AlbumID == 0 {
			changeAlbum = false
		} else {
			changeAlbum = true
		}
	} else {
		if *oldPhoto.AlbumID == *newData.AlbumID {
			changeAlbum = false
		} else {
			changeAlbum = true
		}
	}
	// if (oldPhoto.AlbumID == nil && *newData.AlbumID == 0) || (*oldPhoto.AlbumID == *newData.AlbumID) {
	// 	changeAlbum = false
	// } else {
	// 	changeAlbum = true
	// }

	if *newData.AlbumID == 0 {
		newData.AlbumID = nil
	}

	updateData := map[string]interface{}{
		"title":       newData.Title,
		"description": newData.Description,
		"is_featured": newData.IsFeatured,
		"album_id":    newData.AlbumID, // 当 newData.AlbumID 为 nil 时，这里会正确地将 NULL 写入数据库
		"updated_at":  time.Now(),
	}

	var oldTagsID []int32
	var tagsToAdd []*model.Tag

	newTags := newData.Tags

	if len(newTags) > 0 {
		for _, tag := range newTags {
			tagID, err := findTagIDByName(tag)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					tagsToAdd = append(tagsToAdd, &model.Tag{Name: tag})
				} else {
					return err
				}
			} else {
				oldTagsID = append(oldTagsID, tagID)
			}
		}
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		txq := query.Use(tx)

		txq.Photo.Where(txq.Photo.ID.Eq(newData.ID)).Updates(updateData)
		if changeAlbum && oldPhoto.AlbumID != nil {
			txq.Album.Where(txq.Album.ID.Eq(*oldPhoto.AlbumID)).Updates(map[string]interface{}{"updated_at": time.Now()})
		}
		if changeAlbum && newData.AlbumID != nil {
			txq.Album.Where(txq.Album.ID.Eq(*newData.AlbumID)).Updates(map[string]interface{}{"updated_at": time.Now()})
		}

		if len(tagsToAdd) > 0 {
			if err := txq.Tag.CreateInBatches(tagsToAdd, 100); err != nil {
				return err
			}

			for _, tag := range tagsToAdd {
				oldTagsID = append(oldTagsID, tag.ID)
			}
		}

		if _, err := txq.PhotoTag.Where(txq.PhotoTag.PhotoID.Eq(newData.ID)).Delete(); err != nil {
			return err
		}

		if len(oldTagsID) > 0 {
			var photoTags []*model.PhotoTag
			for _, tagID := range oldTagsID {
				photoTags = append(photoTags, &model.PhotoTag{
					PhotoID: newData.ID,
					TagID:   tagID,
				})
			}
			if err := txq.PhotoTag.CreateInBatches(photoTags, 100); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func deletePhotoByIDStore(photoID int32) error {
	result := database.DB.Model(&model.Photo{}).Where("id = ?", photoID).Delete(&model.Photo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func getPhotosMetaByCursorStore(cursorID int32, limit int, albumID int32) ([]*model.VPhotosWithDetail, error) {
	var photos []*model.VPhotosWithDetail
	db := database.DB.Model(&model.VPhotosWithDetail{})
	if cursorID > 0 {
		var cursorPhoto model.VPhotosWithDetail
		err := database.DB.Model(&model.VPhotosWithDetail{}).Where("id = ?", cursorID).First(&cursorPhoto).Error
		if err != nil {
			return nil, err
		}

		db = db.Where("shot_at < ? OR (shot_at = ? AND id < ?)", cursorPhoto.ShotAt, cursorPhoto.ShotAt, cursorPhoto.ID)
	}

	if albumID != 0 {
		db = db.Where("album_id = ?", albumID)
	}

	err := db.Order("shot_at desc, id desc").Limit(limit).Find(&photos).Error
	if err != nil {
		return nil, err
	}

	return photos, nil

}
