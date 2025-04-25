package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	var config DatabaseConfig

	file, err := os.Open("configs/database.yml")
	if err != nil {
		return nil, fmt.Errorf("❤️ could not open database.yml: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("❤️ could not decode database.yml: %v", err)
	}

	return &config, nil
}