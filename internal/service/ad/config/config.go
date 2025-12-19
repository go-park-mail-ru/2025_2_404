// Package config provides configuration management for the ad service.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig  *PostgresConfig
	AppConfig *AppConfig
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

type AppConfig struct {
	Host        string
	Port        string
	PortAD      string
	PortStorage string
	ImgPath     string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	appCfg := GetAppConfig()

	return &Config{
		DBConfig:  GetPostgresConfig(),
		AppConfig: appCfg,
	}
}

func GetPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DB:       os.Getenv("POSTGRES_DB"),
	}
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		Host:        os.Getenv("APP_HOST"),
		Port:        os.Getenv("APP_PORT"),
		PortAD:      os.Getenv("GRPC_AD_PORT"),
		PortStorage: os.Getenv("GRPC_STORAGE_PORT"),
		ImgPath:     os.Getenv("IMG_PATH"),
	}
}
