package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Config struct {
	Proto    string `default:"tcp"`
	Listen   string `default:"127.0.0.1:1080"`
	LogLevel string `default:"info"`
}

func NewConfig(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return &cfg, err
	}
	err = envconfig.CheckDisallowed(prefix, &cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func PrintUsage(prefix string) {
	var cfg Config
	_ = envconfig.Usage(prefix, &cfg)
	os.Exit(128)
}

func PrintConfig(cfg *Config) {
	j, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(j))
}
