package redis

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis"
)

type redisConf struct {
	rdb *redis.Client
}

func NewRedisStore(c *Config) *redisConf {
	var rc redisConf

	rc.rdb = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	return &rc
}

func (rc *redisConf) KeyExists(key string) (bool, error) {
	exists := rc.rdb.Exists(key)
	if exists.Err() != nil {
		return false, exists.Err()
	}
	return exists.Val() > 0, nil
}

func (rc *redisConf) SetKey(key string, value interface{}, ttl time.Duration) error {
	log.Printf("ttl is %d", ttl)
	result := rc.rdb.Set(key, value, ttl)
	spew.Dump(result)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
