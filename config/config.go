package config

import (
	"os"
)


var (
	JWTSecretKey = getOrDefault("SECRET_KEY", "debug-secret-key")
	DatabaseURL = getOrDefault("DATABASE_URL", "postgresql://localhost:5432?sslmode=disable")
	Debug = !isEnvPresent("CHAT_CONNOR_FUN_PROD")
	Port = getOrDefault("PORT", "4000")
)

func getOrDefault(envVar string, def string) string {
	val, present := os.LookupEnv(envVar)
	if present {
		return val
	}
	return def
}

func isEnvPresent(envVar string) bool {
	_, present := os.LookupEnv(envVar)
	return present
}
