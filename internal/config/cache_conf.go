package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CacheConfig struct {
	Cache struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"cache"`
}

func LoadCacheConfig() (*CacheConfig, error) {
	var config CacheConfig

	file, err := os.Open("configs/cache.yml")
	if err != nil {
		return nil, fmt.Errorf("❤️ Không thể mở file cache.yml: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("❤️ Không thể giải mã file cache.yml: %v", err)
	}

	return &config, nil
}