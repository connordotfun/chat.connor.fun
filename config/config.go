package config

import (
	"os"
)


var (
	JWTSecretKey = getOrDefault("SECRET_KEY", "debug-secret-key")
	DatabaseURL = getOrDefault("DATABASE_URL", "postgresql://user:postgres@localhost:5432?sslmode=disable")
	Debug = !isEnvPresent("CHAT_CONNOR_FUN_PROD")
	Port = getOrDefault("PORT", "4000")
	MailjetPubKey = getOrDefault("MAILJET_PUBLIC_KEY", "0")
	MailjetPrivKey = getOrDefault("MAILJET_PRIVATE_KEY", "0")
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