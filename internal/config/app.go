package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	App struct {
		Name              string `yaml:"name"`
		Host              string `yaml:"host"`
		Port              string `yaml:"port"`
		JWTAccessSecret   string `yaml:"jwt_access_secret"`
		JWTRefreshSecret  string `yaml:"jwt_refresh_secret"`
		ApiPrefix         string `yaml:"api_prefix"`
	} `yaml:"app"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`

	RabbitMQ struct {
		Host string `yaml:"host"`
		Port int `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"rabbitmq"`

	SMTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"smtp"`
}

func LoadAppConfig() (*AppConfig, error) {
	var config AppConfig

	file, err := os.Open("configs/app.yml")
	if err != nil {
		return nil, fmt.Errorf("❤️ Không thể mở file app.yml: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("❤️ Không thể giải mã file app.yml: %v", err)
	}

	return &config, nil
}
