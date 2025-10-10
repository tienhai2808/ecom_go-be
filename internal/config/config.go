package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name        string `yaml:"name"`
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		JWTSecret   string `yaml:"jwt_secret"`
		AccessName  string `yaml:"access_name"`
		RefreshName string `yaml:"refresh_name"`
		ApiPrefix   string `yaml:"api_prefix"`
	} `yaml:"app"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"mysql"`

	Redis struct {
		Addr string `yaml:"addr"`
	} `yaml:"redis"`

	RabbitMQ struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"rabbitmq"`

	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`

	Elasticsearch struct {
		Addresses []string `yaml:"addresses"`
	} `yaml:"elasticsearch"`

	SMTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"smtp"`

	ImageKit struct {
		UrlEndpoint string `yaml:"url_endpoint"`
		PublicKey   string `yaml:"public_key"`
		PrivateKey  string `yaml:"private_key"`
	} `yaml:"imagekit"`

	Cloudinary struct {
		CloudName string `yaml:"cloud_name"`
		ApiKey    string `yaml:"api_key"`
		ApiSecret string `yaml:"api_secret"`
	} `yaml:"cloudinary"`
}

func LoadConfig() (*Config, error) {
	var config Config

	file, err := os.Open("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("không thể mở file config: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("không thể đọc file cấu hình: %w", err)
	}

	return &config, nil
}
