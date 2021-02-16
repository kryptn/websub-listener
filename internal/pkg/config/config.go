package config

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

	// viper.Debug()
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

}

type Listener struct {
	Handler string

	// websub specific
	TopicURL     string
	HubURL       string
	LeaseSeconds int
}

type Emitter struct {
	Handler string

	// slack
	IncomingWebhook string

	// forwarder
	Endpoint string
}

type Store struct {
	// kind should be "redis" or "memory"
	Kind string

	// redis specific settings
	Address  string
	Password string
	DB       int
}

type Config struct {
	PublicHostname string

	Store       Store
	VerifyToken string

	Listeners map[string]Listener `mapstructure:"listener"`
	Emitters  map[string]Emitter  `mapstructure:"emitter"`

	Wires map[string][]string
}

func GetConfig() (*Config, error) {
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}
