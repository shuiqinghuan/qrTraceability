package repository

import (
	"qr-traceability/internal/models"

	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(media *models.MediaFile) error {
	return r.db.Create(media).Error
}

func (r *MediaRepository) GetByID(id uint) (*models.MediaFile, error) {
	var media models.MediaFile
	if err := r.db.First(&media, id).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) ListByBatchID(batchID uint) ([]models.MediaFile, error) {
	var mediaFiles []models.MediaFile
	if err := r.db.Where("batch_id = ?", batchID).Find(&mediaFiles).Error; err != nil {
		return nil, err
	}
	return mediaFiles, nil
}

func (r *MediaRepository) Delete(id uint) error {
	return r.db.Delete(&models.MediaFile{}, id).Error
}
