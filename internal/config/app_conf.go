package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	App struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		JWTAccessSecret string `yaml:"jwt_access_secret"`
		JWTRefreshSecret string `yaml:"jwt_refresh_secret"`
		ApiPrefix string `yaml:"api_prefix"`
		SmtpHost string `yaml:"smtp_host"`
		SmtpPort string `yaml:"smtp_port"`
		SmtpUser string `yaml:"smtp_user"`
		SmtpPass string `yaml:"smtp_pass"`
	} `yaml:"app"`
}

func LoadAppConfig() (*AppConfig, error) {
	var config AppConfig

	file, err := os.Open("configs/app.yml")
	if err != nil {
		return nil, fmt.Errorf("❤️ Không thể mở file app.yml: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("❤️ Không thể giải mã file app.yml: %v", err)
	}

	return &config, nil
}