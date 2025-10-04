package photo

import "time"

type PhotoMeta struct {
	ID           int32
	AlbumID      *int32     // 所属相册ID, 关联albums.id (可以为空, 代表未分类)
	Title        string     // 照片标题
	Description  *string    // 照片描述或背后的故事
	FilePath     string     // 文件存储路径 (例如: /uploads/2024/09/your-photo.jpg)
	FileName     string     // 原始文件名
	FileSize     int32      // 文件大小 (Bytes)
	Width        int32      // 图片宽度 (px)
	Height       int32      // 图片高度 (px)
	IsFeatured   bool       // 是否为精选照片
	ShotAt       *time.Time // 拍摄时间
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CameraID     *int32   // 相机ID, 关联cameras.id
	LensID       *int32   // 镜头ID, 关联lenses.id
	FocalLength  *string  // 焦距 (例如: 85mm)
	Aperture     *string  // 光圈 (例如: f/1.8)
	ShutterSpeed *string  // 快门速度 (例如: 1/1000s)
	Iso          *string  // ISO感光度 (例如: 100)
	ExposureBias *string  // 曝光补偿 (例如: +0.7 EV)
	FlashFired   *bool    // 是否使用闪光灯
	GpsLatitude  *float64 // GPS纬度
	GpsLongitude *float64 // GPS经度
	TagsID       []int32  // 照片标签, 关联tags表 (多对多关系)
}
