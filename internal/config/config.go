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
	Port string
	JwtToken string
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return Config{
		DBConfig:  GetPostgresConfig(),
		AppConfig: GetAppConfig(),
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
		Port: os.Getenv("APP_PORT"),
		JwtToken: os.Getenv("JWT_SECRET"),
	}
}
