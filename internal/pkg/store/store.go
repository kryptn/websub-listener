package store

import (
	"log"
	"time"

	"github.com/kryptn/websub-to-slack/internal/pkg/config"
	"github.com/kryptn/websub-to-slack/internal/pkg/store/memory"
	"github.com/kryptn/websub-to-slack/internal/pkg/store/redis"
)

type Store interface {
	KeyExists(key string) (bool, error)
	SetKey(key string, value interface{}, ttl time.Duration) error
}

func StoreFromConfig(config *config.Config) Store {
	switch config.Store.Kind {
	case "memory":
		return memory.NewMemoryStore()
	case "redis":
		return redis.NewRedisStore(redis.DefaultConfig())
	default:
		log.Fatalf("unable to identify store method %s", config.Store.Kind)
	}
	return nil
}
