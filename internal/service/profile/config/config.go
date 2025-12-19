// Package config provides configuration settings for the profile service.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig      *PostgresConfig
	AppConfig     *AppConfig
	PaymentConfig *PaymentConfig
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

type AppConfig struct {
	Port    string
	ImgPath string
}

type PaymentConfig struct {
	ShopID    string
	SecretKey string
}

func GetConfig() *Config {
	err := godotenv.Load(os.Getenv("ENV_FILE"))
	if err != nil {
		log.Println("Error loading .env file")
	}

	appCfg := GetAppConfig()
	if err != nil {
		log.Println("Error get App configuration")
	}
	return &Config{
		DBConfig:      GetPostgresConfig(),
		AppConfig:     appCfg,
		PaymentConfig: GetPaymentConfig(),
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
		Port:    os.Getenv("GRPC_PORT_PROFILE"),
		ImgPath: os.Getenv("IMG_PATH"),
	}
}

func GetPaymentConfig() *PaymentConfig {
	return &PaymentConfig{
		ShopID:    os.Getenv("YOOKASSA_SHOP_ID"),
		SecretKey: os.Getenv("YOOKASSA_SECRET_KEY"),
	}
}
