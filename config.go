package main

import (
	"fmt"

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

	viper.Debug()
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

}

func getConfig() *Config {
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	listeners := map[string]Subscription{}

	for name, listener := range config.Listeners {
		var l Subscription

		l.Slug = name
		l.TopicURL = listener.TopicURL
		l.VerifyToken = config.VerifyToken
		l.Parser = listener.Parser
		l.Destination = listener.Destination

		if dest, ok := config.Destinations[listener.Destination]; ok {
			l.PostURL = dest
		}

		listeners[name] = l
	}

	config.Listeners = listeners

	return &config

}
