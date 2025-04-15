package utils

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Transfer TransferConfig `yaml:"transfer"`
	Server   ServerConfig   `yaml:"server"`
	AntiBot  AntiBotConfig  `yaml:"antibot"`
	Captcha  CaptchaConfig  `yaml:"captcha"`
	Other    OtherConfig    `yaml:"other"`
}

type TransferConfig struct {
	Compress           bool `yaml:"compress"`
	TimeoutSeconds     int  `yaml:"timeout_seconds"`
	TimeoutIdleSeconds int  `yaml:"timeout_idle_seconds"`
}

type ServerConfig struct {
	UseProfiler  bool   `yaml:"use_profiler"`
	Proxy        string `yaml:"proxy"`
	Origin       string `yaml:"origin"`
	RealIpHeader string `yaml:"real_ip_header"`

	Prefork             bool   `yaml:"prefork"`
	StrictRouting       bool   `yaml:"strict_routing"`
	CaseSensitive       bool   `yaml:"case_sensitive"`
	BodyLimitBytes      int    `yaml:"body_limit_bytes"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
	IdleTimeoutSeconds  int    `yaml:"idle_timeout_seconds"`
	GetIpFromHeader     string `yaml:"get_ip_from_header"`
	EnableKeepAlive     bool   `yaml:"enable_keep_alive"`
}

type AntiBotConfig struct {
	Enabled                     bool `yaml:"enabled"`
	Threshold                   int  `yaml:"threshold"`
	VerificationValidForSeconds int  `yaml:"verification_valid_for_seconds"`
}

type CaptchaConfig struct {
	VerificationValidForSeconds      int `yaml:"verification_valid_for_seconds"`
	MaxFailedAttempts                int `yaml:"max_failed_attempts"`
	MaxFailedAttemptsTimespanSeconds int `yaml:"max_failed_attempts_timespan_seconds"`
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
