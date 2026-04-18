package repository

import (
	"qr-traceability/internal/models"

	"gorm.io/gorm"
)

type HarvestRepository struct {
	db *gorm.DB
}

func NewHarvestRepository(db *gorm.DB) *HarvestRepository {
	return &HarvestRepository{db: db}
}

func (r *HarvestRepository) Create(harvest *models.HarvestQuality) error {
	return r.db.Create(harvest).Error
}

func (r *HarvestRepository) GetByBatchID(batchID uint) (*models.HarvestQuality, error) {
	var harvest models.HarvestQuality
	if err := r.db.Where("batch_id = ?", batchID).First(&harvest).Error; err != nil {
		return nil, err
	}
	return &harvest, nil
}

func (r *HarvestRepository) Update(harvest *models.HarvestQuality) error {
	return r.db.Save(harvest).Error
}

func (r *HarvestRepository) Delete(batchID uint) error {
	return r.db.Where("batch_id = ?", batchID).Delete(&models.HarvestQuality{}).Error
}
