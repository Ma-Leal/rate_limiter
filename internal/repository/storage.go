package repository

import "time"

type Storage interface {
	Incr(key string, window time.Duration) (int64, error)
	Block(key string, duration time.Duration) error
	IsBlocked(key string) (bool, error)
}
