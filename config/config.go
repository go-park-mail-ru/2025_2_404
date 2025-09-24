package config

import (
	"os"
)

func GetPostgresConfig() (string, string, string, string, string) {
	return os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB")
}

func GetAppPort() string {
	return os.Getenv("APP_PORT")
}
