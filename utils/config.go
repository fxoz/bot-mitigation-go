package utils

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	AntiBot AntiBotConfig `yaml:"antibot"`
}

type ServerConfig struct {
	Proxy        string `yaml:"proxy"`
	Origin       string `yaml:"origin"`
	RealIpHeader string `yaml:"real_ip_header"`
}

type AntiBotConfig struct {
	Enabled                  bool `yaml:"enabled"`
	WhitelistDurationSeconds int  `yaml:"whitelist_for_seconds"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig reads the YAML configuration from the specified file.
// The configuration is loaded only once.
func LoadConfig(filePath string) *Config {
	once.Do(func() {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading config file: %v", err)
		}

		tempConfig := &Config{}
		err = yaml.Unmarshal(data, tempConfig)
		if err != nil {
			log.Fatalf("error parsing config file: %v", err)
		}

		config = tempConfig
	})
	return config
}
