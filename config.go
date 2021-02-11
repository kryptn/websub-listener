package main

import (
	"fmt"
	"log"

	"github.com/kryptn/websub-to-slack/internal/pkg/store"
	"github.com/kryptn/websub-to-slack/internal/pkg/store/redis"

	"github.com/spf13/viper"
)

func init() {
	configViper()
}

func configViper() {
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.SetConfigType("toml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/websub/") // path used in docker/k8s
	viper.AddConfigPath(".")            // local path

	// viper.Debug()
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

}

func createSub(name string, listener Subscription, config *Config, s store.Store) Subscription {
	var l Subscription
	l.Slug = name
	log.Printf("setting config slug %s", name)
	l.TopicURL = listener.TopicURL
	l.VerifyToken = config.VerifyToken
	l.Parser = listener.Parser
	l.Destination = listener.Destination
	l.Cache = s

	if dest, ok := config.Destinations[listener.Destination]; ok {
		l.PostURL = dest
	}

	return l
}

func getConfig() *Config {
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	// spew.Dump(config)

	// config.Cache = memory.NewMemoryStore()

	config.Cache = redis.NewRedisStore(redis.DefaultConfig())

	listeners := map[string]Subscription{}

	for name, listener := range config.Listeners {
		listeners[name] = createSub(name, listener, &config, config.Cache)

	}

	config.Listeners = listeners

	return &config

}
