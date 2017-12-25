package config

import (
	"os"
)

type Config struct {
	JWTSecretKey string
	DatabaseURL string
}


func New() Config {
	return Config{os.Getenv("SECRET_KEY"), os.Getenv("DATABASE_URL")}
}