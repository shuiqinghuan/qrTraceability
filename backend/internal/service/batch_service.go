package service

import (
	"fmt"
	"time"

	"qr-traceability/internal/models"
	"qr-traceability/internal/repository"
)

type BatchService struct {
	batchRepo *repository.BatchRepository
	mediaRepo *repository.MediaRepository
	harvestRepo *repository.HarvestRepository
}

func NewBatchService(
	batchRepo *repository.BatchRepository,
	mediaRepo *repository.MediaRepository,
	harvestRepo *repository.HarvestRepository,
) *BatchService {
	return &BatchService{
		batchRepo:   batchRepo,
		mediaRepo:   mediaRepo,
		harvestRepo: harvestRepo,
	}
}

func (s *BatchService) CreateBatch(batch *models.ProductBatch) error {
	// Generate unique ID using product code + timestamp + random string
	// For simplicity, we'll use product code + timestamp
	batch.UniqueID = fmt.Sprintf("%s_%d", batch.BatchNumber, time.Now().UnixNano())
	return s.batchRepo.Create(batch)
}

func (s *BatchService) GetBatchByID(id uint) (*models.ProductBatch, error) {
	return s.batchRepo.GetByID(id)
}

func (s *BatchService) GetBatchByUniqueID(uniqueID string) (*models.ProductBatch, error) {
	return s.batchRepo.GetByUniqueID(uniqueID)
}

func (s *BatchService) ListBatchesByProductID(productID uint, limit, offset int) ([]models.ProductBatch, error) {
	return s.batchRepo.ListByProductID(productID, limit, offset)
}

func (s *BatchService) UpdateBatch(batch *models.ProductBatch) error {
	return s.batchRepo.Update(batch)
}

func (s *BatchService) DeleteBatch(id uint) error {
	return s.batchRepo.Delete(id)
}

func (s *BatchService) AddMedia(batchID uint, media *models.MediaFile) error {
	media.BatchID = batchID
	return s.mediaRepo.Create(media)
}

func (s *BatchService) UpdateHarvestQuality(batchID uint, harvest *models.HarvestQuality) error {
	harvest.BatchID = batchID
	// Check if harvest quality already exists
	existing, err := s.harvestRepo.GetByBatchID(batchID)
	if err == nil && existing != nil {
		// Update existing
		harvest.ID = existing.ID
		return s.harvestRepo.Update(harvest)
	}
	// Create new
	return s.harvestRepo.Create(harvest)
}
