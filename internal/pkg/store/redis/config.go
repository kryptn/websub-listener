package redis

import "time"

// Config struct for the redis memory
type Config struct {
	Addr     string
	Password string
	DB       int

	PostTTL time.Duration
}

//DefaultConfig returns a sane default configuration for Redis
func DefaultConfig() *Config {
	return &Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,

		PostTTL: time.Hour * 6,
	}
}
