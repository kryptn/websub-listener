package memory

import (
	"time"
)

type expiringValue struct {
	value     interface{}
	ttl       time.Duration
	expiresAt time.Time
}

type MemoryStore struct {
	data map[string]expiringValue
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[string]expiringValue{},
	}
}

func (ms *MemoryStore) checkTTLs() {
	var expiredKeys []string

	for key, value := range ms.data {
		if value.expiresAt.IsZero() {
			continue
		}
		if time.Now().After(value.expiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(ms.data, key)
	}
}

func (ms *MemoryStore) KeyExists(key string) (bool, error) {
	ms.checkTTLs()

	_, ok := ms.data[key]
	return ok, nil
}

func (ms *MemoryStore) SetKey(key string, value interface{}, ttl time.Duration) error {
	ms.checkTTLs()

	ev := expiringValue{
		value:     value,
		ttl:       ttl,
		expiresAt: time.Now().Add(ttl),
	}

	ms.data[key] = ev
	return nil
}
