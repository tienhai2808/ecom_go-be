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

	// Đọc file cấu hình
	file, err := os.Open("configs/cache.yml")
	if err != nil {
		return nil, fmt.Errorf("❤️ could not open cache.yml: %v", err)
	}
	defer file.Close()

	// Parse YAML file
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("❤️ could not decode cache.yml: %v", err)
	}

	return &config, nil
}