package repository

import (
	"qr-traceability/internal/models"

	"gorm.io/gorm"
)

type InteractionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

func (r *InteractionRepository) Create(interaction *models.UserInteraction) error {
	return r.db.Create(interaction).Error
}

func (r *InteractionRepository) CountByBatchIDAndAction(batchID uint, actionType string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.UserInteraction{}).Where("batch_id = ? AND action_type = ?", batchID, actionType).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *InteractionRepository) GetInteractionStats(batchID uint) (*models.InteractionStats, error) {
	likeCount, err := r.CountByBatchIDAndAction(batchID, "like")
	if err != nil {
		return nil, err
	}

	shareCount, err := r.CountByBatchIDAndAction(batchID, "share")
	if err != nil {
		return nil, err
	}

	collectCount, err := r.CountByBatchIDAndAction(batchID, "collect")
	if err != nil {
		return nil, err
	}

	return &models.InteractionStats{
		BatchID:      batchID,
		LikeCount:    int(likeCount),
		ShareCount:   int(shareCount),
		CollectCount: int(collectCount),
	}, nil
}
