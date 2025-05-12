package usecase

import (
	"time"

	"github.com/Ma-Leal/rate-limiter/internal/repository"
)

type RateLimiterUseCase struct {
	storage repository.Storage
}

type RateLimiterConfig struct {
	Limit         int
	Window        time.Duration
	BlockDuration time.Duration
}

func NewRateLimiterUseCase(storage repository.Storage) *RateLimiterUseCase {
	return &RateLimiterUseCase{storage: storage}
}

func (u *RateLimiterUseCase) Allow(key string, cfg RateLimiterConfig) (bool, error) {
	blocked, err := u.storage.IsBlocked(key)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	count, err := u.storage.Incr(key, cfg.Window)
	if err != nil {
		return false, err
	}

	if count > int64(cfg.Limit) {
		err := u.storage.Block(key, cfg.BlockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
