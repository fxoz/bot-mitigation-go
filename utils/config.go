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
	Other   OtherConfig   `yaml:"other"`
}

type ServerConfig struct {
	Proxy        string `yaml:"proxy"`
	Origin       string `yaml:"origin"`
	RealIpHeader string `yaml:"real_ip_header"`
}

type AntiBotConfig struct {
	Enabled                     bool `yaml:"enabled"`
	Threshold                   int  `yaml:"threshold"`
	VerificationValidForSeconds int  `yaml:"verification_valid_for_seconds"`
}

type OtherConfig struct {
	ObfuscateJavaScript bool `yaml:"obfuscate_javascript"`
}

var (
	config *Config
	once   sync.Once
)

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
