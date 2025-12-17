// Package config provides configuration management for the auth service.
package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

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
	Host              string
	Port              string
	PortAD            string
	PortStorage       string
	Timeout           time.Duration
	JwtPrivateKeyPath string
	JwtPublicKeyPath  string
	JwtPrivateKey     *ecdsa.PrivateKey
	JwtPublicKey      *ecdsa.PublicKey
}

func GetConfig() *Config {
	err := godotenv.Load(os.Getenv("ENV_FILE"))
	if err != nil {
		log.Println("Error loading .env file")
	}

	appCfg := GetAppConfig()
	if err := LoadJwtPrivateKey(appCfg); err != nil {
		log.Println(fmt.Errorf("ошибка загрузки JWT приватного ключа: %w", err))
	}
	if err := LoadJwtPublicKey(appCfg); err != nil {
		log.Println(fmt.Errorf("ошибка загрузки JWT публичного ключа: %w", err))
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
		Host:        os.Getenv("APP_HOST"),
		Port:        os.Getenv("APP_PORT"),
		PortAD:      os.Getenv("GRPC_AD_PORT"),
		PortStorage: os.Getenv("GRPC_STORAGE_PORT"),
		// Port: os.Getenv("GRPC_PORT_AUTH"),
		JwtPrivateKeyPath: os.Getenv("JWT_PRIVATE_KEY_PATH"),
		JwtPublicKeyPath:  os.Getenv("JWT_PUBLIC_KEY_PATH"),
	}
}

func LoadJwtPrivateKey(cfg *AppConfig) error {
	privKey, err := os.ReadFile(cfg.JwtPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения приватного ключа: %w", err)
	}

	block, _ := pem.Decode(privKey)
	if block == nil {
		return fmt.Errorf("не удалось декодировать PEM блок")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("ошибка парсинга приватного ключа: %w", err)
	}

	cfg.JwtPrivateKey = key

	return nil
}

func LoadJwtPublicKey(cfg *AppConfig) error {
	pubKey, err := os.ReadFile(cfg.JwtPublicKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения публичного ключа: %w", err)
	}

	block, _ := pem.Decode(pubKey)
	if block == nil {
		return fmt.Errorf("не удалось декодировать PEM блок")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("ошибка парсинга публичного ключа: %w", err)
	}

	cfg.JwtPublicKey = key.(*ecdsa.PublicKey)

	return nil
}
