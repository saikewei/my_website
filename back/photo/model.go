package photo

import "time"

type PhotoUpload struct {
	ID           int32      `json:"id"`
	FilePath     string     `json:"file_path"`         // 文件存储路径 (例如: /uploads/2024/09/your-photo.jpg)
	FileName     string     `json:"file_name"`         // 原始文件名
	FileSize     int32      `json:"file_size"`         // 文件大小 (Bytes)
	Width        int32      `json:"width"`             // 图片宽度 (px)
	Height       int32      `json:"height"`            // 图片高度 (px)
	IsFeatured   bool       `json:"is_featured"`       // 是否为精选照片
	ShotAt       *time.Time `json:"shot_at,omitempty"` // 拍摄时间
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Camera       *string    `json:"camera,omitempty"`        // 相机型号 (例如: Canon EOS 5D Mark IV)
	Lens         *string    `json:"lens,omitempty"`          // 镜头型号 (例如: EF 24-70mm f/2.8L II USM)
	FocalLength  *string    `json:"focal_length,omitempty"`  // 焦距 (例如: 85mm)
	Aperture     *string    `json:"aperture,omitempty"`      // 光圈 (例如: f/1.8)
	ShutterSpeed *string    `json:"shutter_speed,omitempty"` // 快门速度 (例如: 1/1000s)
	Iso          *string    `json:"iso,omitempty"`           // ISO感光度 (例如: 100)
	ExposureBias *string    `json:"exposure_bias,omitempty"` // 曝光补偿 (例如: +0.7 EV)
	FlashFired   *bool      `json:"flash_fired,omitempty"`   // 是否使用闪光灯
	GpsLatitude  *float64   `json:"gps_latitude,omitempty"`  // GPS纬度
	GpsLongitude *float64   `json:"gps_longitude,omitempty"` // GPS经度
	Tags         []string   `json:"tags,omitempty"`          // 照片标签, 关联tags表 (多对多关系)
}

type Album struct {
	ID           int32     `json:"id"`
	Title        string    `json:"title"`                    // 相册标题
	Description  *string   `json:"description,omitempty"`    // 相册描述
	CoverPhotoID *int32    `json:"cover_photo_id,omitempty"` // 封面照片ID, 关联photos.id
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PhotoAlbum struct {
	PhotoID int32 `json:"photo_id"`
	AlbumID int32 `json:"album_id"`
}
