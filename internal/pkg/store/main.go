package store

import (
	"time"
)

type Store interface {
	KeyExists(key string) (bool, error)
	SetKey(key string, value interface{}, ttl time.Duration) error
}
