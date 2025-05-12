package tests

import (
	"context"

	"testing"
	"time"

	"github.com/Ma-Leal/rate-limiter/config"
	"github.com/Ma-Leal/rate-limiter/internal/repository"
	"github.com/Ma-Leal/rate-limiter/internal/usecase"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

func setup() (*usecase.RateLimiterUseCase, usecase.RateLimiterConfig) {

	cfg := &config.Config{
		RedisAddr:      "127.0.0.1:6379",
		RateLimitIP:    5,
		RateLimitToken: 10,
		Window:         60 * time.Second,
		BlockDuration:  5 * time.Minute,
	}
	ipCfg := usecase.RateLimiterConfig{
		Limit:         cfg.RateLimitIP,
		Window:        cfg.Window,
		BlockDuration: cfg.BlockDuration,
	}

	storage := repository.NewRedisStorage(cfg.RedisAddr)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(storage)

	return rateLimiterUseCase, ipCfg
}

func TestAllow_RequestWithinLimit(t *testing.T) {
	rateLimiter, ipCfg := setup()

	ipKey := "192.168.1.1"
	allowed, err := rateLimiter.Allow(ipKey, ipCfg)

	assert.NoError(t, err)
	assert.True(t, allowed, "Expected request to be allowed")
}

func TestAllow_RequestExceedsLimit(t *testing.T) {
	rateLimiter, ipCfg := setup()
	ipKey := "192.168.1.1"
	for i := 0; i < 5; i++ {
		_, err := rateLimiter.Allow(ipKey, ipCfg)
		assert.NoError(t, err)
	}
	allowed, err := rateLimiter.Allow(ipKey, ipCfg)

	assert.NoError(t, err)
	assert.False(t, allowed, "Expected request to be blocked")
}

func TestAllow_TokenWithinLimit(t *testing.T) {
	rateLimiter, ipCfg := setup()

	tokenKey := "abc123"
	allowed, err := rateLimiter.Allow(tokenKey, ipCfg)

	assert.NoError(t, err)
	assert.True(t, allowed, "Expected request to be allowed")
}

func TestAllow_TokenExceedsLimit(t *testing.T) {
	rateLimiter, ipCfg := setup()

	tokenKey := "abc123"
	for i := 0; i < 10; i++ {
		_, err := rateLimiter.Allow(tokenKey, ipCfg)
		assert.NoError(t, err)
	}

	allowed, err := rateLimiter.Allow(tokenKey, ipCfg)

	assert.NoError(t, err)
	assert.False(t, allowed, "Expected request to be blocked")
}
