package config

import (
	"log"
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
	Port              string
	ImgPath			  string
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
		Port: os.Getenv("GRPC_PORT_PROFILE"),
		ImgPath: os.Getenv("IMG_PATH"),
	}
}