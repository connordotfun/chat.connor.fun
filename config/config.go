package config

import (
	"os"
)

type Config struct {
	JWTSecretKey string
	DatabaseURL string
	Debug bool
}


func New(debug bool) Config {
	if !debug {
		return Config{os.Getenv("SECRET_KEY"), os.Getenv("DATABASE_URL"), false}
	} else {
		return Config{JWTSecretKey: "hopefully-very-secret-key", DatabaseURL: "todo", Debug: true}
	}
}