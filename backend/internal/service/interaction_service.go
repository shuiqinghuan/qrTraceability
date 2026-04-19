package service

import (
	"context"
	"fmt"
	"time"

	"qr-traceability/internal/models"
	"qr-traceability/internal/repository"

	"github.com/redis/go-redis/v9"
)

type InteractionService struct {
	repo        *repository.InteractionRepository
	redisClient *redis.Client
}

func NewInteractionService(
	repo *repository.InteractionRepository,
	redisClient *redis.Client,
) *InteractionService {
	return &InteractionService{
		repo:        repo,
		redisClient: redisClient,
	}
}

// 防刷配置
const (
	RateLimitWindow = 5 * time.Minute // 时间窗口
	RateLimitCount  = 3                // 窗口内最大操作次数
)

// IsRateLimited 检查是否被限流
func (s *InteractionService) IsRateLimited(ctx context.Context, ip, actionType string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", ip, actionType)

	// 获取当前计数
	count, err := s.redisClient.Get(ctx, key).Int()
	if err == redis.Nil {
		// 第一次操作，设置初始值和过期时间
		if err := s.redisClient.Set(ctx, key, 1, RateLimitWindow).Err(); err != nil {
			return false, err
		}
		return false, nil
	} else if err != nil {
		return false, err
	}

	// 检查是否超过限制
	if count >= RateLimitCount {
		return true, nil
	}

	// 增加计数
	if err := s.redisClient.Incr(ctx, key).Err(); err != nil {
		return false, err
	}

	// 确保键有过期时间
	if err := s.redisClient.Expire(ctx, key, RateLimitWindow).Err(); err != nil {
		return false, err
	}

	return false, nil
}

func (s *InteractionService) RecordInteraction(interaction *models.UserInteraction) error {
	// 检查是否被限流
	ctx := context.Background()
	isLimited, err := s.IsRateLimited(ctx, interaction.IP, interaction.ActionType)
	if err != nil {
		return err
	}
	if isLimited {
		return fmt.Errorf("rate limited: too many requests")
	}

	// 记录交互
	return s.repo.Create(interaction)
}

func (s *InteractionService) GetInteractionStats(batchID uint) (*models.InteractionStats, error) {
	return s.repo.GetInteractionStats(batchID)
}
