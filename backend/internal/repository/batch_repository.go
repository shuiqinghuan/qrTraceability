package repository

import (
	"qr-traceability/internal/models"

	"gorm.io/gorm"
)

type BatchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

func (r *BatchRepository) Create(batch *models.ProductBatch) error {
	return r.db.Create(batch).Error
}

func (r *BatchRepository) GetByID(id uint) (*models.ProductBatch, error) {
	var batch models.ProductBatch
	if err := r.db.Preload("Product").Preload("MediaFiles").Preload("HarvestQuality").First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *BatchRepository) GetByUniqueID(uniqueID string) (*models.ProductBatch, error) {
	var batch models.ProductBatch
	if err := r.db.Preload("Product").Preload("MediaFiles").Preload("HarvestQuality").Where("unique_id = ?", uniqueID).First(&batch).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *BatchRepository) ListByProductID(productID uint, limit, offset int) ([]models.ProductBatch, error) {
	var batches []models.ProductBatch
	if err := r.db.Where("product_id = ?", productID).Limit(limit).Offset(offset).Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *BatchRepository) Update(batch *models.ProductBatch) error {
	return r.db.Save(batch).Error
}

func (r *BatchRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductBatch{}, id).Error
}
