package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client  *redis.Client
	context context.Context
}

func NewRedisStorage(addr string) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStorage{
		client:  rdb,
		context: context.Background(),
	}
}

func (r *RedisStorage) Incr(key string, window time.Duration) (int64, error) {
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(r.context, key)
	pipe.Expire(r.context, key, window)
	_, err := pipe.Exec(r.context)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

func (r *RedisStorage) Block(key string, duration time.Duration) error {
	blockKey := key + ":block"
	return r.client.Set(r.context, blockKey, "1", duration).Err()
}

func (r *RedisStorage) IsBlocked(key string) (bool, error) {
	blockKey := key + ":block"
	val, err := r.client.Get(r.context, blockKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
