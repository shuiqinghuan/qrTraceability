package models

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品基本信息
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Code        string `gorm:"uniqueIndex;not null" json:"code"` // 产品编码，如4395
	Name        string `gorm:"not null" json:"name"`             // 产品名称，如枣甜5号
	Description string `json:"description"`                        // 产品描述
	
	Batches []ProductBatch `gorm:"foreignKey:ProductID" json:"batches,omitempty"`
}

// ProductBatch 产品批次信息
type ProductBatch struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	ProductID   uint   `gorm:"not null" json:"product_id"`
	BatchNumber string `gorm:"not null" json:"batch_number"` // 批次号
	UniqueID    string `gorm:"uniqueIndex;not null" json:"unique_id"` // 唯一标识，用于生成二维码
	
	PlantingLocation string    `json:"planting_location"` // 定植地点
	PlantingDate     time.Time `json:"planting_date"`     // 定植时间
	
	Product      Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	MediaFiles   []MediaFile    `gorm:"foreignKey:BatchID" json:"media_files,omitempty"`
	HarvestQuality *HarvestQuality `gorm:"foreignKey:BatchID" json:"harvest_quality,omitempty"`
	Interactions []UserInteraction `gorm:"foreignKey:BatchID" json:"interactions,omitempty"`
}

// MediaFile 媒体文件信息
type MediaFile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	BatchID   uint   `gorm:"not null" json:"batch_id"`
	Type      string `gorm:"not null" json:"type"` // image或video
	URL       string `gorm:"not null" json:"url"`  // 文件存储URL
	Thumbnail string `json:"thumbnail"`             // 缩略图URL（仅图片）
	
	Batch ProductBatch `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
}

// HarvestQuality 采收质量信息
type HarvestQuality struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	BatchID           uint      `gorm:"uniqueIndex;not null" json:"batch_id"`
	HarvestStartDate  time.Time `json:"harvest_start_date"`  // 采收起始时间
	HarvestEndDate    time.Time `json:"harvest_end_date"`    // 采收终止时间
	SugarContent      float64   `json:"sugar_content"`      // 糖度
	Weight            float64   `json:"weight"`            // 重量
	Taste             string    `json:"taste"`             // 口感
	SuitableFor       string    `json:"suitable_for"`       // 适应人群
	QualitySummary    string    `json:"quality_summary"`    // 品质小结
	
	Batch ProductBatch `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
}

// UserInteraction 用户交互数据
type UserInteraction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	BatchID    uint   `gorm:"not null" json:"batch_id"`
	IP         string `gorm:"not null" json:"ip"`         // 用户IP
	DeviceID   string `json:"device_id"`                   // 设备标识
	ActionType string `gorm:"not null" json:"action_type"` // like, share, collect
	
	Batch ProductBatch `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
}

// InteractionStats 交互统计
type InteractionStats struct {
	BatchID    uint `json:"batch_id"`
	LikeCount  int  `json:"like_count"`
	ShareCount int  `json:"share_count"`
	CollectCount int `json:"collect_count"`
}
